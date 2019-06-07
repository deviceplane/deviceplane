package service

import (
	"github.com/Knetic/govaluate"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/pkg/errors"
)

type deviceLabelParameters []models.DeviceLabel

func (p deviceLabelParameters) Get(name string) (interface{}, error) {
	for _, deviceLabel := range []models.DeviceLabel(p) {
		if deviceLabel.Key == name {
			return deviceLabel.Value, nil
		}
	}
	return "", nil
}

type emptyParameters struct{}

func (p emptyParameters) Get(name string) (interface{}, error) {
	return "", nil
}

func validateSchedulingRule(schedulingRule string) error {
	expression, err := govaluate.NewEvaluableExpression(schedulingRule)
	if err != nil {
		return errors.Wrap(err, "invalid scheduling rule")
	}

	result, err := expression.Eval(emptyParameters{})
	if err != nil {
		return errors.Wrap(err, "evaluate scheduling rule")
	}

	if _, ok := result.(bool); result == nil || !ok {
		return errors.New("scheduling rule should evaluate to a boolean value")
	}

	return nil
}
