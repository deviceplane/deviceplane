package query

import (
	"encoding/base64"
	"encoding/json"
)

func FilterDevices(devicesMap []map[string]interface{}, query map[string][]string) []map[string]interface{} {
	filters := filtersFromQuery(query)

	if len(filters) == 0 {
		return devicesMap
	}

	filteredDevices := make([]map[string]interface{}, 0, len(devicesMap))

	for _, device := range devicesMap {
		valid := true
		for _, filter := range filters {
			matchesFilter := false
			for _, condition := range filter {
				if matchesFilter {
					break
				}
				property := condition["property"].(string)
				operator := condition["operator"].(string)
				if property == "label" {
					key := condition["key"].(string)

					switch operator {
					case "is":
						value := condition["value"].(string)
						for _, label := range device["labels"].([]map[string]interface{}) {
							if label["key"] == key && label["value"] == value {
								matchesFilter = true
								break
							}
						}
					case "is not":
						value := condition["value"].(string)
						found := false
						for _, label := range device["labels"].([]map[string]interface{}) {
							if label["key"] == key && label["value"] == value {
								found = true
								break
							}
						}
						if !found {
							matchesFilter = true
						}
					case "key is":
						for _, label := range device["labels"].([]map[string]interface{}) {
							if label["key"] == key {
								matchesFilter = true
								break
							}
						}
					case "key is not":
						found := false
						for _, label := range device["labels"].([]map[string]interface{}) {
							if label["key"] == key {
								found = true
								break
							}
						}
						if !found {
							matchesFilter = true
						}
					}
				} else {
					value := condition["value"].(string)
					switch operator {
					case "is":
						matchesFilter = device[property] == value
					case "is not":
						matchesFilter = device[property] != value
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

func filtersFromQuery(query map[string][]string) [][]map[string]interface{} {
	var filters [][]map[string]interface{} = nil
	for key, values := range query {
		if key == "filter" {
			for _, encodedFilter := range values {
				var filter []map[string]interface{}
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
