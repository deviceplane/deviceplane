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
)

type Query []Filter

type Filter []Condition

type Condition struct {
	Type   ConditionType          `json:"type"`
	Params map[string]interface{} `json:"params"`
}

type ConditionType string

const (
	DevicePropertyCondition = ConditionType("DevicePropertyCondition")
	LabelValueCondition     = ConditionType("LabelValueCondition")
	LabelExistenceCondition = ConditionType("LabelExistenceCondition")
)

type DevicePropertyConditionParams struct {
	Property string   `json:"property"`
	Operator Operator `json:"operator"`
	Value    string   `json:"value"`
}

type LabelValueConditionParams struct {
	Key      string   `json:"key"`
	Operator Operator `json:"operator"`
	Value    string   `json:"value"`
}

type LabelExistenceConditionParams struct {
	Key      string   `json:"key"`
	Operator Operator `json:"operator"`
}

type Operator string

const (
	OperatorIs    = Operator("is")
	OperatorIsNot = Operator("is not")

	OperatorExists    = Operator("exists")
	OperatorNotExists = Operator("does not exist")
)

func FilterDevices(devices []models.Device, query Query) ([]models.Device, error) {
	filteredDevices := make([]models.Device, 0)
	for _, device := range devices {
		match, err := DeviceMatchesQuery(device, query)
		if err != nil {
			return nil, err
		}
		if match {
			filteredDevices = append(filteredDevices, device)
		}
	}
	return filteredDevices, nil
}

func DeviceMatchesQuery(device models.Device, query Query) (bool, error) {
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

func deviceMatchesFilter(device models.Device, filter Filter) (bool, error) {
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

func deviceMatchesCondition(device models.Device, condition Condition) (bool, error) {
	switch condition.Type {
	case DevicePropertyCondition:
		var params DevicePropertyConditionParams
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
		case OperatorIs:
			return match, nil
		case OperatorIsNot:
			return !match, nil
		}
		return false, ErrOperatorNotSupported

	case LabelValueCondition:
		var params LabelValueConditionParams
		err := utils.JSONConvert(condition.Params, &params)
		if err != nil {
			return false, err
		}

		value, ok := device.Labels[params.Key]
		valueMatches := bool(ok && value == params.Value)

		switch params.Operator {
		case OperatorIs:
			return valueMatches, nil
		case OperatorIsNot:
			return !valueMatches, nil
		}
		return false, ErrOperatorNotSupported

	case LabelExistenceCondition:
		var params LabelExistenceConditionParams
		err := utils.JSONConvert(condition.Params, &params)
		if err != nil {
			return false, err
		}

		_, ok := device.Labels[params.Key]
		switch params.Operator {
		case OperatorExists:
			return ok, nil
		case OperatorNotExists:
			return !ok, nil
		}
		return false, ErrOperatorNotSupported
	}
	return false, ErrConditionNotSupported
}

func FiltersFromQuery(query map[string][]string) ([]Filter, error) {
	var filters []Filter

	for key, values := range query {
		if key == "filter" {
			for _, encodedFilter := range values {
				bytes, err := base64.StdEncoding.DecodeString(encodedFilter)
				if err != nil {
					return nil, err
				}

				var filter Filter
				if err := json.Unmarshal(bytes, &filter); err != nil {
					return nil, err
				}

				filters = append(filters, filter)
			}
		}
	}

	return filters, nil
}
