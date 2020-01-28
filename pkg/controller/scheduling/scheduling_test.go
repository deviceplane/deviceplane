package scheduling

import (
	"testing"

	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/stretchr/testify/require"
)

type Scenario struct {
	in             []models.Device
	schedulingRule models.SchedulingRule
	out            []models.ScheduledDevice
}

func testScenario(t *testing.T, scenario Scenario) {
	t.Helper()
	scheduledDevices, err := GetScheduledDevices(scenario.in, scenario.schedulingRule)
	require.NoError(t, err)
	require.Equal(t, scenario.out, scheduledDevices)
}

func TestScheduleAllDevices(t *testing.T) {
	testScenario(t, Scenario{
		in: []models.Device{
			models.Device{
				ID: "one",
				Labels: map[string]string{
					"a": "b",
				},
			},
			models.Device{
				ID: "two",
				Labels: map[string]string{
					"a": "b",
				},
			},
		},
		schedulingRule: models.SchedulingRule{
			ScheduleType:     models.ScheduleTypeAllDevices,
			DefaultReleaseID: "1",
		},
		out: []models.ScheduledDevice{
			models.ScheduledDevice{
				Device: models.Device{
					ID: "one",
					Labels: map[string]string{
						"a": "b",
					},
				},
				ReleaseID: "1",
			},
			models.ScheduledDevice{
				Device: models.Device{
					ID: "two",
					Labels: map[string]string{
						"a": "b",
					},
				},
				ReleaseID: "1",
			},
		},
	})
}

func TestScheduleAllDevicesWithReleaseSelectors(t *testing.T) {
	testScenario(t, Scenario{
		in: []models.Device{
			models.Device{
				ID: "one",
				Labels: map[string]string{
					"a": "b",
				},
			},
			models.Device{
				ID: "two",
				Labels: map[string]string{
					"a": "b",
				},
			},
			models.Device{
				ID: "three",
				Labels: map[string]string{
					"a": "c",
				},
			},
			models.Device{
				ID: "four",
				Labels: map[string]string{
					"a": "d",
				},
			},
		},
		schedulingRule: models.SchedulingRule{
			ScheduleType: models.ScheduleTypeAllDevices,
			ReleaseSelectors: []models.ReleaseSelector{
				models.ReleaseSelector{
					Query: models.Query{
						models.Filter{
							models.Condition{
								Type: models.LabelValueCondition,
								Params: map[string]interface{}{
									"key":      "a",
									"operator": models.OperatorIs,
									"value":    "b",
								},
							},
						},
					},
					ReleaseID: "pinned",
				},
				models.ReleaseSelector{
					Query: models.Query{
						models.Filter{
							models.Condition{
								Type: models.LabelValueCondition,
								Params: map[string]interface{}{
									"key":      "a",
									"operator": models.OperatorIs,
									"value":    "c",
								},
							},
						},
					},
					ReleaseID: "canary",
				},
			},
			DefaultReleaseID: "1",
		},
		out: []models.ScheduledDevice{
			models.ScheduledDevice{
				Device: models.Device{
					ID: "one",
					Labels: map[string]string{
						"a": "b",
					},
				},
				ReleaseID: "pinned",
			},
			models.ScheduledDevice{
				Device: models.Device{
					ID: "two",
					Labels: map[string]string{
						"a": "b",
					},
				},
				ReleaseID: "pinned",
			},
			models.ScheduledDevice{
				Device: models.Device{
					ID: "three",
					Labels: map[string]string{
						"a": "c",
					},
				},
				ReleaseID: "canary",
			},
			models.ScheduledDevice{
				Device: models.Device{
					ID: "four",
					Labels: map[string]string{
						"a": "d",
					},
				},
				ReleaseID: "1",
			},
		},
	})
}

func TestScheduleNoDevices(t *testing.T) {
	testScenario(t, Scenario{
		in: []models.Device{
			models.Device{
				ID: "one",
				Labels: map[string]string{
					"a": "b",
				},
			},
			models.Device{
				ID: "two",
				Labels: map[string]string{
					"a": "b",
				},
			},
		},
		schedulingRule: models.SchedulingRule{
			ScheduleType:     models.ScheduleTypeNoDevices,
			DefaultReleaseID: "1",
		},
		out: []models.ScheduledDevice{},
	})
}

func TestScheduleWithQuery(t *testing.T) {
	testScenario(t, Scenario{
		in: []models.Device{
			models.Device{
				ID:     "one",
				Status: models.DeviceStatusOnline,
				Labels: map[string]string{
					"a": "b",
				},
			},
			models.Device{
				ID:     "two",
				Status: models.DeviceStatusOnline,
				Labels: map[string]string{
					"a": "a",
				},
			},
		},
		schedulingRule: models.SchedulingRule{
			ScheduleType: models.ScheduleTypeConditional,
			ConditionalQuery: &models.Query{
				models.Filter{
					models.Condition{
						Type: models.LabelValueCondition,
						Params: map[string]interface{}{
							"key":      "a",
							"operator": models.OperatorIs,
							"value":    "b",
						},
					},
				},
			},
			ReleaseSelectors: nil,
			DefaultReleaseID: "1",
		},
		out: []models.ScheduledDevice{
			models.ScheduledDevice{
				Device: models.Device{
					ID:     "one",
					Status: models.DeviceStatusOnline,
					Labels: map[string]string{
						"a": "b",
					},
				},
				ReleaseID: "1",
			},
		},
	})
}

func TestScheduleWithSimplePinnedQuery(t *testing.T) {
	testScenario(t, Scenario{
		in: []models.Device{
			models.Device{
				ID:     "one",
				Status: models.DeviceStatusOnline,
				Labels: map[string]string{
					"a": "b",
				},
			},
			models.Device{
				ID:     "two",
				Status: models.DeviceStatusOnline,
				Labels: map[string]string{
					"a": "a",
				},
			},
		},
		schedulingRule: models.SchedulingRule{
			ScheduleType: models.ScheduleTypeConditional,
			ConditionalQuery: &models.Query{
				models.Filter{
					models.Condition{
						Type: models.LabelValueCondition,
						Params: map[string]interface{}{
							"key":      "a",
							"operator": models.OperatorIs,
							"value":    "b",
						},
					},
				},
			},
			ReleaseSelectors: []models.ReleaseSelector{
				models.ReleaseSelector{
					Query: models.Query{
						models.Filter{
							models.Condition{
								Type: models.LabelValueCondition,
								Params: map[string]interface{}{
									"key":      "a",
									"operator": models.OperatorIs,
									"value":    "b",
								},
							},
						},
					},
					ReleaseID: "pinned",
				},
			},
			DefaultReleaseID: "1",
		},
		out: []models.ScheduledDevice{
			models.ScheduledDevice{
				Device: models.Device{
					ID:     "one",
					Status: models.DeviceStatusOnline,
					Labels: map[string]string{
						"a": "b",
					},
				},
				ReleaseID: "pinned",
			},
		},
	})
}

func TestScheduleWithComplexPinnedQuery(t *testing.T) {
	testScenario(t, Scenario{
		in: []models.Device{
			models.Device{
				ID:     "one",
				Status: models.DeviceStatusOnline,
				Labels: map[string]string{
					"a": "b",
				},
			},
			models.Device{
				ID:     "two",
				Status: models.DeviceStatusOnline,
				Labels: map[string]string{
					"a": "c",
				},
			},
			models.Device{
				ID:     "three",
				Status: models.DeviceStatusOnline,
				Labels: map[string]string{
					"a": "d",
				},
			},
			models.Device{
				ID:     "four",
				Status: models.DeviceStatusOnline,
				Labels: map[string]string{
					"a": "test",
				},
			},
		},
		schedulingRule: models.SchedulingRule{
			ScheduleType: models.ScheduleTypeConditional,
			ConditionalQuery: &models.Query{
				models.Filter{
					models.Condition{
						Type: models.LabelValueCondition,
						Params: map[string]interface{}{
							"key":      "a",
							"operator": models.OperatorIs,
							"value":    "b",
						},
					},
					models.Condition{
						Type: models.LabelValueCondition,
						Params: map[string]interface{}{
							"key":      "a",
							"operator": models.OperatorIs,
							"value":    "c",
						},
					},
					models.Condition{
						Type: models.LabelValueCondition,
						Params: map[string]interface{}{
							"key":      "a",
							"operator": models.OperatorIs,
							"value":    "d",
						},
					},
				},
			},
			ReleaseSelectors: []models.ReleaseSelector{
				models.ReleaseSelector{
					Query: models.Query{
						models.Filter{
							models.Condition{
								Type: models.LabelValueCondition,
								Params: map[string]interface{}{
									"key":      "a",
									"operator": models.OperatorIs,
									"value":    "b",
								},
							},
							models.Condition{
								Type: models.LabelValueCondition,
								Params: map[string]interface{}{
									"key":      "a",
									"operator": models.OperatorIs,
									"value":    "f",
								},
							},
						},
					},
					ReleaseID: "pinned",
				},
				models.ReleaseSelector{
					Query: models.Query{
						models.Filter{
							models.Condition{
								Type: models.LabelValueCondition,
								Params: map[string]interface{}{
									"key":      "a",
									"operator": models.OperatorIs,
									"value":    "c",
								},
							},
						},
					},
					ReleaseID: "pinned-two",
				},
				models.ReleaseSelector{
					Query: models.Query{
						models.Filter{
							models.Condition{
								Type: models.LabelValueCondition,
								Params: map[string]interface{}{
									"key":      "a",
									"operator": models.OperatorIs,
									"value":    "test",
								},
							},
						},
					},
					ReleaseID: "pinned-three",
				},
			},
			DefaultReleaseID: "1",
		},
		out: []models.ScheduledDevice{
			models.ScheduledDevice{
				Device: models.Device{
					ID:     "one",
					Status: models.DeviceStatusOnline,
					Labels: map[string]string{
						"a": "b",
					},
				},
				ReleaseID: "pinned",
			},
			models.ScheduledDevice{
				Device: models.Device{
					ID:     "two",
					Status: models.DeviceStatusOnline,
					Labels: map[string]string{
						"a": "c",
					},
				},
				ReleaseID: "pinned-two",
			},
			models.ScheduledDevice{
				Device: models.Device{
					ID:     "three",
					Status: models.DeviceStatusOnline,
					Labels: map[string]string{
						"a": "d",
					},
				},
				ReleaseID: "1",
			},
		},
	})
}

func TestIndividualDeviceScheduling(t *testing.T) {
	schedulingRule := models.SchedulingRule{
		ScheduleType: models.ScheduleTypeConditional,
		ConditionalQuery: &models.Query{
			models.Filter{
				models.Condition{
					Type: models.LabelValueCondition,
					Params: map[string]interface{}{
						"key":      "a",
						"operator": models.OperatorIs,
						"value":    "b",
					},
				},
				models.Condition{
					Type: models.LabelValueCondition,
					Params: map[string]interface{}{
						"key":      "a",
						"operator": models.OperatorIs,
						"value":    "c",
					},
				},
				models.Condition{
					Type: models.LabelValueCondition,
					Params: map[string]interface{}{
						"key":      "a",
						"operator": models.OperatorIs,
						"value":    "d",
					},
				},
			},
		},
		ReleaseSelectors: []models.ReleaseSelector{
			models.ReleaseSelector{
				Query: models.Query{
					models.Filter{
						models.Condition{
							Type: models.LabelValueCondition,
							Params: map[string]interface{}{
								"key":      "a",
								"operator": models.OperatorIs,
								"value":    "b",
							},
						},
						models.Condition{
							Type: models.LabelValueCondition,
							Params: map[string]interface{}{
								"key":      "a",
								"operator": models.OperatorIs,
								"value":    "f",
							},
						},
					},
				},
				ReleaseID: "pinned",
			},
			models.ReleaseSelector{
				Query: models.Query{
					models.Filter{
						models.Condition{
							Type: models.LabelValueCondition,
							Params: map[string]interface{}{
								"key":      "a",
								"operator": models.OperatorIs,
								"value":    "c",
							},
						},
					},
				},
				ReleaseID: "pinned-two",
			},
			models.ReleaseSelector{
				Query: models.Query{
					models.Filter{
						models.Condition{
							Type: models.LabelValueCondition,
							Params: map[string]interface{}{
								"key":      "a",
								"operator": models.OperatorIs,
								"value":    "test",
							},
						},
					},
				},
				ReleaseID: "pinned-three",
			},
		},
		DefaultReleaseID: "1",
	}

	testIndividualScheduling := func(t *testing.T, d models.Device, sr models.SchedulingRule, is bool, sd *models.ScheduledDevice, err error) {
		isSched, scheduledDevice, schedErr := IsApplicationScheduled(d, sr)
		require.Equal(t, is, isSched, "is scheduled")
		if sd != nil && scheduledDevice != nil {
			require.Equal(t, *sd, *scheduledDevice, "scheduled device")
		} else {
			require.Equal(t, sd, scheduledDevice, "scheduled device")
		}
		require.Equal(t, err, schedErr, "error when scheduling")
	}

	testIndividualScheduling(
		t,
		models.Device{
			ID:     "one",
			Status: models.DeviceStatusOnline,
			Labels: map[string]string{
				"a": "b",
			},
		},
		schedulingRule,

		true,
		&models.ScheduledDevice{
			Device: models.Device{
				ID:     "one",
				Status: models.DeviceStatusOnline,
				Labels: map[string]string{
					"a": "b",
				},
			},
			ReleaseID: "pinned",
		},
		nil,
	)

	testIndividualScheduling(
		t,
		models.Device{
			ID:     "two",
			Status: models.DeviceStatusOnline,
			Labels: map[string]string{
				"a": "c",
			},
		},
		schedulingRule,

		true,
		&models.ScheduledDevice{
			Device: models.Device{
				ID:     "two",
				Status: models.DeviceStatusOnline,
				Labels: map[string]string{
					"a": "c",
				},
			},
			ReleaseID: "pinned-two",
		},
		nil,
	)

	testIndividualScheduling(
		t,
		models.Device{
			ID:     "three",
			Status: models.DeviceStatusOnline,
			Labels: map[string]string{
				"a": "d",
			},
		},
		schedulingRule,

		true,
		&models.ScheduledDevice{
			Device: models.Device{
				ID:     "three",
				Status: models.DeviceStatusOnline,
				Labels: map[string]string{
					"a": "d",
				},
			},
			ReleaseID: "1",
		},
		nil,
	)

	testIndividualScheduling(
		t,
		models.Device{
			ID:     "four",
			Status: models.DeviceStatusOnline,
			Labels: map[string]string{
				"a": "test",
			},
		},
		schedulingRule,

		false,
		nil,
		nil,
	)
}
