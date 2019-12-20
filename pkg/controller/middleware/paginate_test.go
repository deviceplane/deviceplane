package middleware

import (
	"testing"

	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/stretchr/testify/require"
)

type PaginationScenario struct {
	desc       string
	after      string
	paginateOn string
	pageSize   int
	in         []interface{}
	out        []interface{}
	shouldErr  bool
}

func testPaginationScenario(t *testing.T, s PaginationScenario) {
	t.Helper()
	out, err := paginateAfter(s.after, s.paginateOn, s.pageSize, s.in)
	if s.shouldErr {
		require.Error(t, err, s.desc)
	} else {
		require.NoError(t, err, s.desc)
		require.Equal(t, s.out, out, s.desc)
	}
}

func TestPagination(t *testing.T) {
	scenarios := []PaginationScenario{
		PaginationScenario{
			desc:       "Test getting first page (empty)",
			after:      "",
			paginateOn: "id",
			pageSize:   1,
			shouldErr:  false,
			in: []interface{}{
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
				models.Device{
					ID:     "device_d",
					Status: models.DeviceStatusOffline,
				},
			},
			out: []interface{}{
				models.Device{
					ID:     "device_a",
					Status: models.DeviceStatusOffline,
				},
			},
		},
		PaginationScenario{
			desc:       "Test getting second page",
			after:      "device_a",
			paginateOn: "id",
			pageSize:   1,
			shouldErr:  false,
			in: []interface{}{
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
				models.Device{
					ID:     "device_d",
					Status: models.DeviceStatusOffline,
				},
			},
			out: []interface{}{
				models.Device{
					ID:     "device_b",
					Status: models.DeviceStatusOffline,
				},
			},
		},
		PaginationScenario{
			desc:       "Test getting last page",
			after:      "device_c",
			paginateOn: "id",
			pageSize:   1,
			shouldErr:  false,
			in: []interface{}{
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
				models.Device{
					ID:     "device_d",
					Status: models.DeviceStatusOffline,
				},
			},
			out: []interface{}{
				models.Device{
					ID:     "device_d",
					Status: models.DeviceStatusOffline,
				},
			},
		},
		PaginationScenario{
			desc:       "Test getting last + 1 page (nonexistent)",
			after:      "device_d",
			paginateOn: "id",
			pageSize:   1,
			shouldErr:  false,
			in: []interface{}{
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
				models.Device{
					ID:     "device_d",
					Status: models.DeviceStatusOffline,
				},
			},
			out: []interface{}{},
		},
		PaginationScenario{
			desc:       "Test getting (nonexistent)",
			after:      "device_yeet",
			paginateOn: "id",
			pageSize:   1,
			shouldErr:  true,
			in: []interface{}{
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
				models.Device{
					ID:     "device_d",
					Status: models.DeviceStatusOffline,
				},
			},
		},
		PaginationScenario{
			desc:       "Test page size of 3",
			after:      "",
			paginateOn: "id",
			pageSize:   3,
			shouldErr:  false,
			in: []interface{}{
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
				models.Device{
					ID:     "device_d",
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
		PaginationScenario{
			desc:       "Test getting a partial page",
			after:      "device_c",
			paginateOn: "id",
			pageSize:   3,
			shouldErr:  false,
			in: []interface{}{
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
				models.Device{
					ID:     "device_d",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_e",
					Status: models.DeviceStatusOffline,
				},
			},
			out: []interface{}{
				models.Device{
					ID:     "device_d",
					Status: models.DeviceStatusOffline,
				},
				models.Device{
					ID:     "device_e",
					Status: models.DeviceStatusOffline,
				},
			},
		},
		PaginationScenario{
			desc:       "Paginate on nonexistent value on basic type",
			after:      "",
			paginateOn: "id",
			pageSize:   1,
			shouldErr:  true,
			in: []interface{}{
				1,
				2,
				3,
			},
		},
		PaginationScenario{
			desc:       "Paginate on nonexistent value on struct",
			after:      "",
			paginateOn: "not-a-json-tag",
			pageSize:   1,
			shouldErr:  true,
			in: []interface{}{
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
	}

	for _, s := range scenarios {
		t.Run(s.desc, func(t *testing.T) {
			testPaginationScenario(t, s)
		})
	}
}
