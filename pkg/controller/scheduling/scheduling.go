package scheduling

import (
	"encoding/base64"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/deviceplane/deviceplane/pkg/controller/query"
	"github.com/deviceplane/deviceplane/pkg/models"
)

var (
	ErrInvalidScheduleType   = errors.New("invalid schedule type")
	ErrInvalidConditionValue = errors.New("invalid condition value")
	ErrInvalidReleaseID      = errors.New("invalid release ID")

	ErrNonexistentSchedulingRule = errors.New("nonexistent scheduling rule")
)

func SchedulingRuleFromQuery(query map[string][]string) (*models.SchedulingRule, error) {
	schedulingRuleStr, exists := query["schedulingRule"]
	if !exists || len(schedulingRuleStr) == 0 {
		return nil, nil
	}

	bytes, err := base64.StdEncoding.DecodeString(schedulingRuleStr[0])
	if err != nil {
		return nil, err
	}

	var schedulingRule models.SchedulingRule
	if err := json.Unmarshal(bytes, &schedulingRule); err != nil {
		return nil, err
	}

	return &schedulingRule, nil
}

func IsApplicationScheduled(device models.Device, schedulingRule models.SchedulingRule) (bool, *models.ScheduledDevice, error) {
	scheduledDevices, err := GetScheduledDevices([]models.Device{device}, schedulingRule)
	if err != nil {
		return false, nil, err
	}
	if len(scheduledDevices) == 0 {
		return false, nil, nil
	}
	return true, &scheduledDevices[0], nil
}

func GetScheduledDevices(devices []models.Device, schedulingRule models.SchedulingRule) ([]models.ScheduledDevice, error) {
	var selectedDevices []models.Device

	switch schedulingRule.ScheduleType {
	case models.ScheduleTypeNoDevices:
		return []models.ScheduledDevice{}, nil

	case models.ScheduleTypeAllDevices:
		selectedDevices = devices

	case models.ScheduleTypeConditional:
		if schedulingRule.ConditionalQuery == nil {
			return nil, ErrInvalidConditionValue
		}

		var err error
		selectedDevices, _, err = query.QueryDevices(devices, *schedulingRule.ConditionalQuery)
		if err != nil {
			return nil, errors.Wrap(err, "filtering by schedule query")
		}

	default:
		return nil, ErrInvalidScheduleType
	}

	if len(selectedDevices) == 0 {
		return []models.ScheduledDevice{}, nil
	}

	var scheduledDevices []models.ScheduledDevice

	// Go through release selectors
	for _, releaseSelector := range schedulingRule.ReleaseSelectors {
		releasePinnedDevices, newSelectedDevices, err := query.QueryDevices(selectedDevices, releaseSelector.Query)
		if err != nil {
			return nil, errors.Wrap(err, "filtering by release query")
		}

		for _, pinnedDevice := range releasePinnedDevices {
			scheduledDevices = append(scheduledDevices, models.ScheduledDevice{
				Device:    pinnedDevice,
				ReleaseID: releaseSelector.ReleaseID,
			})
		}

		selectedDevices = newSelectedDevices
	}

	for _, defaultReleaseDevice := range selectedDevices {
		scheduledDevices = append(scheduledDevices, models.ScheduledDevice{
			Device:    defaultReleaseDevice,
			ReleaseID: schedulingRule.DefaultReleaseID,
		})
	}

	return scheduledDevices, nil
}

func ValidateSchedulingRule(schedulingRule models.SchedulingRule, releaseIdExists func(string) (bool, error)) (
	validationErr error,
	err error,
) {
	switch schedulingRule.ScheduleType {
	case models.ScheduleTypeNoDevices:
		break

	case models.ScheduleTypeAllDevices:
		break

	case models.ScheduleTypeConditional:
		if schedulingRule.ConditionalQuery == nil {
			return ErrInvalidConditionValue, nil
		}
		err := query.ValidateQuery(*schedulingRule.ConditionalQuery)
		if err != nil {
			return errors.Wrap(err, "filtering by schedule query"), nil
		}
	default:
		return nil, ErrInvalidScheduleType
	}

	if schedulingRule.DefaultReleaseID != models.LatestRelease {
		var exists bool
		exists, err = releaseIdExists(schedulingRule.DefaultReleaseID)
		if err != nil {
			return nil, err
		}

		if !exists {
			return errors.Wrapf(ErrInvalidReleaseID, "default release %s", schedulingRule.DefaultReleaseID), nil
		}
	}

	// Go through release selectors
	for _, releaseSelector := range schedulingRule.ReleaseSelectors {
		err := query.ValidateQuery(releaseSelector.Query)
		if err != nil {
			return errors.Wrap(err, "filtering by release query"), nil
		}

		if releaseSelector.ReleaseID != models.LatestRelease {
			var exists bool
			exists, err = releaseIdExists(releaseSelector.ReleaseID)
			if err != nil {
				return nil, err
			}

			if !exists {
				return errors.Wrapf(ErrInvalidReleaseID, "release %s", schedulingRule.DefaultReleaseID), nil
			}
		}
	}

	return nil, nil
}
