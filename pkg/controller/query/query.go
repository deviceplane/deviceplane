package query

import (
	"encoding/base64"
	"encoding/json"
)

type Operator string

const (
	OperatorIs             = Operator("is")
	OperatorIsNot          = Operator("is not")
	OperatorHasKey         = Operator("has key")
	OperatorDoesNotHaveKey = Operator("does not have key")
)

type Filter []Condition

type Condition struct {
	Property string   `json:"property"`
	Operator Operator `json:"operator"`
	Key      string   `json:"key"`
	Value    string   `json:"value"`
}

func FilterDevices(devicesMap []map[string]interface{}, filters []Filter) []map[string]interface{} {
	filteredDevices := make([]map[string]interface{}, 0)

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
						for key, value := range device["labels"].(map[string]interface{}) {
							if key == condition.Key && value == condition.Value {
								matchesFilter = true
								break
							}
						}
						for key, value := range device["labels"].(map[string]interface{}) {
							if key == condition.Key && value == condition.Value {
								matchesFilter = true
								break
							}
						}
					case OperatorIsNot:
						found := false
						for key, value := range device["labels"].(map[string]interface{}) {
							if key == condition.Key && value == condition.Value {
								found = true
								break
							}
						}
						if !found {
							matchesFilter = true
						}
					case OperatorHasKey:
						for key := range device["labels"].(map[string]interface{}) {
							if key == condition.Key {
								matchesFilter = true
								break
							}
						}
					case OperatorDoesNotHaveKey:
						found := false
						for key := range device["labels"].(map[string]interface{}) {
							if key == condition.Key {
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

	return filteredDevices
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
