package query

import (
	"encoding/base64"
	"encoding/json"

	"github.com/deviceplane/deviceplane/pkg/models"
)

type Operator string

type FilterCondition struct {
	Property string
	Operator Operator
	Key      string
	Value    string
}

const (
	OperatorIs       = Operator("is")
	OperatorIsNot    = Operator("is not")
	OperatorKeyIs    = Operator("key is")
	OperatorKeyIsNot = Operator("key is not")
)

func FilterDevices(devices []models.DeviceWithLabels, filters [][]FilterCondition) ([]map[string]interface{}, bool) {
	var devicesMap []map[string]interface{}

	jsonBytes, err := json.Marshal(devices)

	if err != nil {
		return nil, true
	}

	json.Unmarshal(jsonBytes, &devicesMap)

	filteredDevices := make([]map[string]interface{}, 0, len(devicesMap))

	for _, device := range devicesMap {
		valid := true
		for _, filter := range filters {
			matchesFilter := false
			for _, condition := range filter {
				if matchesFilter {
					break
				}
				if condition.Property == "label" {
					switch condition.Operator {
					case OperatorIs:
						for _, label := range device["labels"].([]map[string]interface{}) {
							if label["key"] == condition.Key && label["value"] == condition.Value {
								matchesFilter = true
								break
							}
						}
					case OperatorIsNot:
						found := false
						for _, label := range device["labels"].([]map[string]interface{}) {
							if label["key"] == condition.Key && label["value"] == condition.Value {
								found = true
								break
							}
						}
						if !found {
							matchesFilter = true
						}
					case OperatorKeyIs:
						for _, label := range device["labels"].([]map[string]interface{}) {
							if label["key"] == condition.Key {
								matchesFilter = true
								break
							}
						}
					case OperatorKeyIsNot:
						found := false
						for _, label := range device["labels"].([]map[string]interface{}) {
							if label["key"] == condition.Key {
								found = true
								break
							}
						}
						if !found {
							matchesFilter = true
						}
					}
				} else {
					switch condition.Operator {
					case OperatorIs:
						matchesFilter = device[condition.Property] == condition.Value
					case OperatorIsNot:
						matchesFilter = device[condition.Property] != condition.Value
					}
				}
			}
			if !matchesFilter {
				valid = false
				break
			}
		}

		if valid {
			filteredDevices = append(filteredDevices, device)
		}
	}

	return filteredDevices, false
}

func FiltersFromQuery(query map[string][]string) [][]FilterCondition {
	var filters [][]FilterCondition
	for key, values := range query {
		if key == "filter" {
			for _, encodedFilter := range values {
				var filter []FilterCondition
				bytes, err := base64.StdEncoding.DecodeString(encodedFilter)
				if err == nil {
					err := json.Unmarshal(bytes, &filter)
					if err == nil {
						filters = append(filters, filter)
					}
				}
			}
		}
	}

	return filters
}
