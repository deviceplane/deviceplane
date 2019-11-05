package middleware

import (
	"testing"

	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/stretchr/testify/require"
)

type OrderScenario struct {
	desc      string
	orderBy   string
	direction orderDirection
	in        []interface{}
	out       []interface{}
}

func testOrderScenario(t *testing.T, s OrderScenario) {
	t.Helper()
	err := order(s.orderBy, s.direction, s.in)
	require.NoError(t, err, s.desc)
	require.Equal(t, s.out, s.in, s.desc)
}

func expectErrInOrderScenario(t *testing.T, s OrderScenario) {
	t.Helper()
	err := order(s.orderBy, s.direction, s.in)
	require.Error(t, err, s.desc)
}

func TestOrdering(t *testing.T) {
	exampleStr := "asdf"
	exampleStrTwo := "bsdf"
	scenarios := []OrderScenario{
		OrderScenario{
			desc:      "Test ascending ordering on string ID",
			direction: OrderAscending,
			orderBy:   "id",
			in: []interface{}{
				models.Device{
					ID:     "device_a",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_c",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_b",
					Status: models.DeviceStatusOffline,
				},
			},
			out: []interface{}{
				models.Device{
					ID:     "device_a",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_b",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_c",
					Status: models.DeviceStatusOffline,
				},
			},
		},
		OrderScenario{
			desc:      "Test descending ordering on string ID",
			direction: OrderDescending,
			orderBy:   "id",
			in: []interface{}{
				models.Device{
					ID:     "device_a",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_c",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_b",
					Status: models.DeviceStatusOffline,
				},
			},
			out: []interface{}{
				models.Device{
					ID:     "device_c",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_b",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_a",
					Status: models.DeviceStatusOffline,
				},
			},
		},
		OrderScenario{
			desc:      "Order on name",
			direction: OrderAscending,
			orderBy:   "name",
			in: []interface{}{
				models.Device{
					ID:     "device_a",
					Name:   "aaa",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_c",
					Name:   "ccc",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_b",
					Name:   "bbb",
					Status: models.DeviceStatusOffline,
				},
			},
			out: []interface{}{
				models.Device{
					ID:     "device_a",
					Name:   "aaa",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_b",
					Name:   "bbb",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_c",
					Name:   "ccc",
					Status: models.DeviceStatusOffline,
				},
			},
		},
		OrderScenario{
			desc:      "Reverse order on name",
			direction: OrderAscending,
			orderBy:   "name",
			in: []interface{}{
				models.Device{
					ID:     "device_a",
					Name:   "aaa",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_c",
					Name:   "ccc",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_b",
					Name:   "bbb",
					Status: models.DeviceStatusOffline,
				},
			},
			out: []interface{}{
				models.Device{
					ID:     "device_a",
					Name:   "aaa",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_b",
					Name:   "bbb",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_c",
					Name:   "ccc",
					Status: models.DeviceStatusOffline,
				},
			},
		},
		OrderScenario{
			desc:      "Order on type equivalent to string",
			direction: OrderAscending,
			orderBy:   "status",
			in: []interface{}{
				models.Device{
					ID:     "device_a",
					Name:   "aaa",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_c",
					Name:   "ccc",
					Status: models.DeviceStatusOnline,
				},
				models.Device{
					ID:     "device_b",
					Name:   "bbb",
					Status: models.DeviceStatusOffline,
				},
			},
			out: []interface{}{
				models.Device{
					ID:     "device_a",
					Name:   "aaa",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_b",
					Name:   "bbb",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_c",
					Name:   "ccc",
					Status: models.DeviceStatusOnline,
				},
			},
		},
		OrderScenario{
			desc:      "Order on possibly-nil value",
			direction: OrderAscending,
			orderBy:   "registrationTokenId",
			in: []interface{}{
				models.Device{
					ID:                  "device_a",
					Name:                "aaa",
					RegistrationTokenID: &exampleStr,
					Status:              models.DeviceStatusOffline,
				},
				models.Device{
					ID:                  "device_c",
					Name:                "ccc",
					RegistrationTokenID: &exampleStrTwo,
					Status:              models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_b",
					Name:   "bbb",
					Status: models.DeviceStatusOffline,
				},
			},
			out: []interface{}{
				models.Device{
					ID:     "device_b",
					Name:   "bbb",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:                  "device_a",
					Name:                "aaa",
					RegistrationTokenID: &exampleStr,
					Status:              models.DeviceStatusOffline,
				},
				models.Device{
					ID:                  "device_c",
					Name:                "ccc",
					RegistrationTokenID: &exampleStrTwo,
					Status:              models.DeviceStatusOffline,
				},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.desc, func(t *testing.T) {
			testOrderScenario(t, s)
		})
	}
}

func TestErrorsInOrdering(t *testing.T) {
	scenarios := []OrderScenario{
		OrderScenario{
			desc:      "Order by missing field",
			direction: OrderAscending,
			orderBy:   "asdfeeeeep",
			in: []interface{}{
				models.Device{
					ID:     "device_a",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_c",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_b",
					Status: models.DeviceStatusOffline,
				},
			},
		},
		OrderScenario{
			desc:      "Order on field with unsupported data type",
			direction: OrderAscending,
			orderBy:   "info",
			in: []interface{}{
				models.Device{
					ID:     "device_a",
					Name:   "aaa",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_c",
					Name:   "ccc",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_b",
					Name:   "bbb",
					Status: models.DeviceStatusOffline,
				},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.desc, func(t *testing.T) {
			expectErrInOrderScenario(t, s)
		})
	}
}
