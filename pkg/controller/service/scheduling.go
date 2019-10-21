package service

import (
	"github.com/Knetic/govaluate"
	"github.com/pkg/errors"
)

type deviceLabelParameters map[string]string

func (p deviceLabelParameters) Get(name string) (interface{}, error) {
	label, ok := p[name]
	if !ok {
		return "", nil
	}

	return label, nil
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
