package query

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilterDevices(t *testing.T) {
	t.Run("standard", func(t *testing.T) {
		for _, scenario := range []struct {
			in      []map[string]interface{}
			filters []Filter
			out     []map[string]interface{}
		}{
			{
				in: []map[string]interface{}{
					{
						"status": "online",
					},
				},
				out: []map[string]interface{}{
					{
						"status": "online",
					},
				},
			},
			{
				in: []map[string]interface{}{
					{
						"status": "online",
					},
				},
				filters: []Filter{
					Filter([]Condition{
						{
							Property: "status",
							Operator: OperatorIs,
							Value:    "online",
						},
					}),
				},
				out: []map[string]interface{}{
					{
						"status": "online",
					},
				},
			},
			{
				in: []map[string]interface{}{
					{
						"status": "offline",
					},
				},
				filters: []Filter{
					Filter([]Condition{
						{
							Property: "status",
							Operator: OperatorIs,
							Value:    "online",
						},
					}),
				},
				out: []map[string]interface{}{},
			},
		} {
			require.Equal(t, scenario.out, FilterDevices(scenario.in, scenario.filters))
		}
	})

	t.Run("label", func(t *testing.T) {
		for _, scenario := range []struct {
			in      []map[string]interface{}
			filters []Filter
			out     []map[string]interface{}
		}{
			{
				in: []map[string]interface{}{
					{
						"labels": map[string]interface{}{
							"a": "b",
						},
					},
				},
				filters: []Filter{
					Filter([]Condition{
						{
							Property: "label",
							Operator: OperatorIs,
							Key:      "a",
							Value:    "b",
						},
					}),
				},
				out: []map[string]interface{}{
					{
						"labels": map[string]interface{}{
							"a": "b",
						},
					},
				},
			},
			{
				in: []map[string]interface{}{
					{
						"labels": map[string]interface{}{
							"a": "b",
						},
					},
				},
				filters: []Filter{
					Filter([]Condition{
						{
							Property: "label",
							Operator: OperatorHasKey,
							Key:      "a",
						},
					}),
				},
				out: []map[string]interface{}{
					{
						"labels": map[string]interface{}{
							"a": "b",
						},
					},
				},
			},
			{
				in: []map[string]interface{}{
					{
						"labels": map[string]interface{}{
							"a": "b",
						},
					},
				},
				filters: []Filter{
					Filter([]Condition{
						{
							Property: "label",
							Operator: OperatorDoesNotHaveKey,
							Key:      "c",
						},
					}),
				},
				out: []map[string]interface{}{
					{
						"labels": map[string]interface{}{
							"a": "b",
						},
					},
				},
			},
			{
				in: []map[string]interface{}{
					{
						"labels": map[string]interface{}{
							"a": "b",
						},
					},
				},
				filters: []Filter{
					Filter([]Condition{
						{
							Property: "label",
							Operator: OperatorIs,
							Key:      "a",
							Value:    "x",
						},
					}),
				},
				out: []map[string]interface{}{},
			},
			{
				in: []map[string]interface{}{
					{
						"labels": map[string]interface{}{
							"a": "b",
						},
					},
				},
				filters: []Filter{
					Filter([]Condition{
						{
							Property: "label",
							Operator: OperatorHasKey,
							Key:      "c",
						},
					}),
				},
				out: []map[string]interface{}{},
			},
			{
				in: []map[string]interface{}{
					{
						"labels": map[string]interface{}{
							"a": "b",
						},
					},
				},
				filters: []Filter{
					Filter([]Condition{
						{
							Property: "label",
							Operator: OperatorDoesNotHaveKey,
							Key:      "a",
						},
					}),
				},
				out: []map[string]interface{}{},
			},
		} {
			require.Equal(t, scenario.out, FilterDevices(scenario.in, scenario.filters))
		}
	})
}

func TestFiltersFromQuery(t *testing.T) {
	filtersA := Filter([]Condition{
		{
			Property: "status",
			Operator: OperatorIs,
			Value:    "online",
		},
		{
			Property: "status",
			Operator: OperatorIs,
			Value:    "offline",
		},
	})

	filtersB := Filter([]Condition{
		{
			Property: "status",
			Operator: OperatorIs,
			Value:    "online",
		},
	})

	jsonFilterA, _ := json.Marshal(filtersA)
	encodedFilterA := base64.StdEncoding.EncodeToString(jsonFilterA)
	jsonFilterB, _ := json.Marshal(filtersB)
	encodedFilterB := base64.StdEncoding.EncodeToString(jsonFilterB)

	query := map[string][]string{
		"filter": []string{
			encodedFilterA,
			encodedFilterB,
		},
	}

	result, err := FiltersFromQuery(query)
	require.NoError(t, err)

	require.Len(t, result, 2)
	require.Len(t, result[0], 2)
	require.Len(t, result[1], 1)
	require.Equal(t, filtersA[0], result[0][0])
	require.Equal(t, filtersA[1], result[0][1])
	require.Equal(t, filtersB[0], result[1][0])
}
