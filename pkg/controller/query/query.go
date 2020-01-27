package query

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/utils"
)

var (
	ErrConditionNotSupported = errors.New("condition not supported")
	ErrOperatorNotSupported  = errors.New("operator not supported")
	ErrPropertyNotSupported  = errors.New("device property not supported")
	ErrNoEmptyFields         = errors.New("fields should not be empty")
)

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

func QueryDevices(devices []models.Device, query models.Query) (selectedDevices []models.Device, unselectedDevices []models.Device, err error) {
	selectedDevices = make([]models.Device, 0)
	unselectedDevices = make([]models.Device, 0)

	for _, device := range devices {
		match, err := DeviceMatchesQuery(device, query)
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

func DeviceMatchesQuery(device models.Device, query models.Query) (bool, error) {
	for _, filter := range query {
		match, err := deviceMatchesFilter(device, filter)
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}
	return true, nil
}

func deviceMatchesFilter(device models.Device, filter models.Filter) (bool, error) {
	for _, condition := range filter {
		match, err := deviceMatchesCondition(device, condition)
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
		return ErrOperatorNotSupported

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
		return ErrOperatorNotSupported

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
		return ErrOperatorNotSupported
	}
	return ErrConditionNotSupported
}

func deviceMatchesCondition(device models.Device, condition models.Condition) (bool, error) {
	switch condition.Type {
	case models.DevicePropertyCondition:
		var params models.DevicePropertyConditionParams
		err := utils.JSONConvert(condition.Params, &params)
		if err != nil {
			return false, err
		}

		// We can either turn device into a map[string]interface, or make
		// a switch statement with cases for each property.
		// I'm going with the former just because it's easier.
		var deviceMap map[string]interface{}
		err = utils.JSONConvert(device, &deviceMap)
		if err != nil {
			return false, err
		}

		value, exists := deviceMap[params.Property]
		if !exists {
			return false, ErrPropertyNotSupported
		}

		match := value == params.Value
		switch params.Operator {
		case models.OperatorIs:
			return match, nil
		case models.OperatorIsNot:
			return !match, nil
		}
		return false, ErrOperatorNotSupported

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
		return false, ErrOperatorNotSupported

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
		return false, ErrOperatorNotSupported
	}
	return false, ErrConditionNotSupported
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
