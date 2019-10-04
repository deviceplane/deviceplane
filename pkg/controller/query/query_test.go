package query

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func MockDevices() [][]map[string]interface{} {
	var devicesA []map[string]interface{}
	var devicesB []map[string]interface{}
	var devicesC []map[string]interface{}

	devicesA = append(devicesA, map[string]interface{}{
		"id": "A",
	}, map[string]interface{}{
		"id": "B",
	},
		map[string]interface{}{
			"id": "C",
		},
		map[string]interface{}{
			"id": "E",
		})

	devicesB = append(devicesB, map[string]interface{}{
		"id": "A",
	}, map[string]interface{}{
		"id": "B",
	})

	devicesC = append(devicesC, map[string]interface{}{
		"id": "B",
	}, map[string]interface{}{
		"id": "E",
	},
		map[string]interface{}{
			"id": "A",
		},
		map[string]interface{}{
			"id": "G",
		},
		map[string]interface{}{
			"id": "H",
		})

	devices := [][]map[string]interface{}{devicesA, devicesB, devicesC}

	return devices
}

func MockFilters() []map[string]interface{} {
	filters := []map[string]interface{}{{"property": "status",
		"operator": "is",
		"value":    "online"},
		{"property": "status",
			"operator": "is",
			"value":    "offline"}}

	return filters
}

func TestFilterDevices(t *testing.T) {

}

func TestFiltersFromQuery(t *testing.T) {
	filters := MockFilters()

	filtersA := []map[string]interface{}{
		filters[0], filters[1],
	}

	filtersB := []map[string]interface{}{
		filters[1],
	}

	jsonFilterA, _ := json.Marshal(filtersA)
	encodedFilterA := base64.StdEncoding.EncodeToString(jsonFilterA)

	jsonFilterB, _ := json.Marshal(filtersB)
	encodedFilterB := base64.StdEncoding.EncodeToString(jsonFilterB)

	query := map[string][]string{
		"filter": []string{encodedFilterA, encodedFilterB},
	}

	result := filtersFromQuery(query)

	require.True(t, len(result) == 2)
	require.True(t, len(result[0]) == 2)
	require.True(t, len(result[1]) == 1)
	require.True(t, reflect.DeepEqual(result[0][0], filters[0]))
	require.True(t, reflect.DeepEqual(result[0][1], filters[1]))
	require.True(t, reflect.DeepEqual(result[1][0], filters[1]))
}
