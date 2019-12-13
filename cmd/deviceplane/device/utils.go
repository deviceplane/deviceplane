package device

import (
	"fmt"
	"strings"

	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/deviceplane/deviceplane/pkg/utils"
)

func parseTextFilter(text string) (models.Filter, error) {
	if strings.HasPrefix(text, "labels.") {
		text = text[len("labels."):]

		if i := strings.Index(text, "!="); i != -1 {
			condition := models.Condition{
				Type: models.LabelValueCondition,
			}

			err := utils.JSONConvert(models.LabelValueConditionParams{
				Key:      text[:i],
				Operator: models.OperatorIsNot,
				Value:    text[i+len("!="):],
			},
				&condition.Params,
			)
			if err != nil {
				return nil, err
			}

			return []models.Condition{condition}, nil
		}
		if i := strings.Index(text, "="); i != -1 {
			condition := models.Condition{
				Type: models.LabelValueCondition,
			}

			err := utils.JSONConvert(models.LabelValueConditionParams{
				Key:      text[:i],
				Operator: models.OperatorIs,
				Value:    text[i+len("="):],
			},
				&condition.Params,
			)
			if err != nil {
				return nil, err
			}

			return []models.Condition{condition}, nil
		}
	}
	if i := strings.Index(text, "!="); i != -1 {
		condition := models.Condition{
			Type: models.DevicePropertyCondition,
		}

		err := utils.JSONConvert(models.DevicePropertyConditionParams{
			Property: text[:i],
			Operator: models.OperatorIsNot,
			Value:    text[i+len("!="):],
		},
			&condition.Params,
		)
		if err != nil {
			return nil, err
		}

		return []models.Condition{condition}, nil
	}
	if i := strings.Index(text, "="); i != -1 {
		condition := models.Condition{
			Type: models.DevicePropertyCondition,
		}

		err := utils.JSONConvert(models.DevicePropertyConditionParams{
			Property: text[:i],
			Operator: models.OperatorIs,
			Value:    text[i+len("="):],
		},
			&condition.Params,
		)
		if err != nil {
			return nil, err
		}

		return []models.Condition{condition}, nil
	}

	return nil, fmt.Errorf(`invalid or missing operator in filter "%s"`, text)
}
