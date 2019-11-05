package middleware

import (
	"testing"

	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/stretchr/testify/require"
)

type PaginationScenario struct {
	desc        string
	pageNumber  int
	pageSize    int
	in          []interface{}
	out         []interface{}
	shouldExist bool
}

func testPaginationScenario(t *testing.T, s PaginationScenario) {
	t.Helper()
	out, exists := paginate(s.pageNumber, s.pageSize, s.in)
	if s.shouldExist {
		require.True(t, exists)
		require.Equal(t, s.out, out)
	} else {
		require.False(t, exists, out)
	}
}

func TestPagination(t *testing.T) {
	scenarios := []PaginationScenario{
		PaginationScenario{
			desc:        "Test getting first page",
			pageNumber:  0,
			pageSize:    1,
			shouldExist: true,
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
			desc:        "Test getting second page",
			pageNumber:  1,
			pageSize:    1,
			shouldExist: true,
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
			desc:        "Test getting last page",
			pageNumber:  3,
			pageSize:    1,
			shouldExist: true,
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
			desc:        "Test getting last + 1 page (nonexistent)",
			pageNumber:  4,
			pageSize:    1,
			shouldExist: false,
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
			desc:        "Test getting -1 page (nonexistent)",
			pageNumber:  -1,
			pageSize:    1,
			shouldExist: false,
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
			desc:        "Test page size of 3",
			pageNumber:  0,
			pageSize:    3,
			shouldExist: true,
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
			desc:        "Test getting a partial page",
			pageNumber:  1,
			pageSize:    3,
			shouldExist: true,
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
	}

	for _, s := range scenarios {
		t.Run(s.desc, func(t *testing.T) {
			testPaginationScenario(t, s)
		})
	}
}
