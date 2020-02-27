package query

import (
	"encoding/base64"
	"encoding/json"

	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/utils"
	"github.com/pkg/errors"
)

var (
	ErrConditionInvalid    = errors.New("invalid condition")
	ErrOperatorInvalid     = errors.New("invalid operator")
	ErrPropertyInvalid     = errors.New("invalid device property")
	ErrServiceStateInvalid = errors.New("invalid service state")

	ErrNoEmptyFields = errors.New("fields should not be empty")
)

type QueryDependencies struct {
	DeviceApplicationStatuses map[string]map[string]models.DeviceApplicationStatus
	DeviceServiceStates       map[string]map[string]map[string]models.DeviceServiceState
}

func ValidateQuery(query models.Query) error {
	for _, filter := range query {
		for _, condition := range filter {
			err := validateCondition(condition)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func QueryDevices(deps QueryDependencies, devices []models.Device, query models.Query) (selectedDevices []models.Device, unselectedDevices []models.Device, err error) {
	selectedDevices = make([]models.Device, 0)
	unselectedDevices = make([]models.Device, 0)

	for _, device := range devices {
		match, err := DeviceMatchesQuery(deps, device, query)
		if err != nil {
			return nil, nil, err
		}

		if match {
			selectedDevices = append(selectedDevices, device)
		} else {
			unselectedDevices = append(unselectedDevices, device)
		}
	}
	return selectedDevices, unselectedDevices, nil
}

func DeviceMatchesQuery(deps QueryDependencies, device models.Device, query models.Query) (bool, error) {
	for _, filter := range query {
		match, err := deviceMatchesFilter(deps, device, filter)
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}
	return true, nil
}

func deviceMatchesFilter(deps QueryDependencies, device models.Device, filter models.Filter) (bool, error) {
	for _, condition := range filter {
		match, err := deviceMatchesCondition(deps, device, condition)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}

	return false, nil
}

func validateCondition(condition models.Condition) error {
	switch condition.Type {
	case models.DevicePropertyCondition:
		var params models.DevicePropertyConditionParams
		err := utils.JSONConvert(condition.Params, &params)
		if err != nil {
			return err
		}

		if params.Property == "" {
			return ErrNoEmptyFields
		}
		if params.Value == "" {
			return ErrNoEmptyFields
		}

		switch params.Operator {
		case models.OperatorIs:
			return nil
		case models.OperatorIsNot:
			return nil
		}
		return ErrOperatorInvalid

	case models.LabelValueCondition:
		var params models.LabelValueConditionParams
		err := utils.JSONConvert(condition.Params, &params)
		if err != nil {
			return err
		}

		if params.Key == "" {
			return ErrNoEmptyFields
		}
		if params.Value == "" {
			return ErrNoEmptyFields
		}

		switch params.Operator {
		case models.OperatorIs:
			return nil
		case models.OperatorIsNot:
			return nil
		}
		return ErrOperatorInvalid

	case models.LabelExistenceCondition:
		var params models.LabelExistenceConditionParams
		err := utils.JSONConvert(condition.Params, &params)
		if err != nil {
			return err
		}

		if params.Key == "" {
			return ErrNoEmptyFields
		}

		switch params.Operator {
		case models.OperatorExists:
			return nil
		case models.OperatorNotExists:
			return nil
		}
		return ErrOperatorInvalid

	case models.ApplicationReleaseCondition:
		var params models.ApplicationReleaseConditionParams
		err := utils.JSONConvert(condition.Params, &params)
		if err != nil {
			return err
		}

		if params.ApplicationID == "" || params.ReleaseID == "" {
			return ErrNoEmptyFields
		}

		switch params.Operator {
		case models.OperatorIs:
			return nil
		case models.OperatorIsNot:
			return nil
		}
		return ErrOperatorInvalid

	case models.ServiceStateCondition:
		var params models.ServiceStateConditionParams
		err := utils.JSONConvert(condition.Params, &params)
		if err != nil {
			return err
		}

		if params.ApplicationID == "" || params.Service == "" {
			return ErrNoEmptyFields
		}

		switch params.Operator {
		case models.OperatorIs:
		case models.OperatorIsNot:
		default:
			return ErrOperatorInvalid
		}

		if !models.AllServiceStates[params.ServiceState] {
			return ErrServiceStateInvalid
		}
		return nil
	}
	return ErrConditionInvalid
}

func deviceMatchesCondition(deps QueryDependencies, device models.Device, condition models.Condition) (bool, error) {
	switch condition.Type {
	case models.DevicePropertyCondition:
		var params models.DevicePropertyConditionParams
		err := utils.JSONConvert(condition.Params, &params)
		if err != nil {
			return false, err
		}

		var deviceMap map[string]interface{}
		err = utils.JSONConvert(device, &deviceMap)
		if err != nil {
			return false, err
		}

		value, exists := deviceMap[params.Property]
		if !exists {
			return false, ErrPropertyInvalid
		}

		match := value == params.Value
		switch params.Operator {
		case models.OperatorIs:
			return match, nil
		case models.OperatorIsNot:
			return !match, nil
		}
		return false, ErrOperatorInvalid

	case models.LabelValueCondition:
		var params models.LabelValueConditionParams
		err := utils.JSONConvert(condition.Params, &params)
		if err != nil {
			return false, err
		}

		value, ok := device.Labels[params.Key]
		valueMatches := bool(ok && value == params.Value)

		switch params.Operator {
		case models.OperatorIs:
			return valueMatches, nil
		case models.OperatorIsNot:
			return !valueMatches, nil
		}
		return false, ErrOperatorInvalid

	case models.LabelExistenceCondition:
		var params models.LabelExistenceConditionParams
		err := utils.JSONConvert(condition.Params, &params)
		if err != nil {
			return false, err
		}

		_, ok := device.Labels[params.Key]
		switch params.Operator {
		case models.OperatorExists:
			return ok, nil
		case models.OperatorNotExists:
			return !ok, nil
		}
		return false, ErrOperatorInvalid

	case models.ApplicationReleaseCondition:
		var params models.ApplicationReleaseConditionParams
		err := utils.JSONConvert(condition.Params, &params)
		if err != nil {
			return false, err
		}

		if params.ApplicationID == "" || params.ReleaseID == "" {
			return false, ErrNoEmptyFields
		}

		applicationStatus, exists := deps.DeviceApplicationStatuses[device.ID][params.ApplicationID]
		existsFunc := func() (bool, error) {
			if !exists {
				return false, nil
			}
			if params.ReleaseID == models.AnyApplicationRelease {
				return true, nil
			}
			return params.ReleaseID == applicationStatus.CurrentReleaseID, nil
		}

		switch params.Operator {
		case models.OperatorIs:
			exists, err := existsFunc()
			return exists, err
		case models.OperatorIsNot:
			exists, err := existsFunc()
			return !exists, err
		}
		return false, ErrOperatorInvalid

	case models.ServiceStateCondition:
		var params models.ServiceStateConditionParams
		err := utils.JSONConvert(condition.Params, &params)
		if err != nil {
			return false, err
		}

		if params.ApplicationID == "" || params.Service == "" {
			return false, ErrNoEmptyFields
		}

		exists := true
		deviceServiceState, exists := deps.DeviceServiceStates[device.ID][params.ApplicationID][params.Service]

		switch params.Operator {
		case models.OperatorIs:
			if !exists {
				return false, nil
			}
			return deviceServiceState.State == params.ServiceState, nil
		case models.OperatorIsNot:
			if !exists {
				return true, nil
			}
			return deviceServiceState.State != params.ServiceState, nil
		}
		return false, ErrOperatorInvalid
	}
	return false, ErrConditionInvalid
}

func FiltersFromQuery(query map[string][]string) ([]models.Filter, error) {
	var filters []models.Filter

	for key, values := range query {
		if key == "filter" {
			for _, encodedFilter := range values {
				bytes, err := base64.StdEncoding.DecodeString(encodedFilter)
				if err != nil {
					return nil, err
				}

				var filter models.Filter
				if err := json.Unmarshal(bytes, &filter); err != nil {
					return nil, err
				}

				filters = append(filters, filter)
			}
		}
	}

	return filters, nil
}
