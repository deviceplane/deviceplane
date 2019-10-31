package models

type Query []Filter

type Filter []Condition

type Condition struct {
	Type   ConditionType          `json:"type"`
	Params map[string]interface{} `json:"params"`
}

type ConditionType string

const (
	DevicePropertyCondition = ConditionType("DevicePropertyCondition")
	LabelValueCondition     = ConditionType("LabelValueCondition")
	LabelExistenceCondition = ConditionType("LabelExistenceCondition")
)

type DevicePropertyConditionParams struct {
	Property string   `json:"property"`
	Operator Operator `json:"operator"`
	Value    string   `json:"value"`
}

type LabelValueConditionParams struct {
	Key      string   `json:"key"`
	Operator Operator `json:"operator"`
	Value    string   `json:"value"`
}

type LabelExistenceConditionParams struct {
	Key      string   `json:"key"`
	Operator Operator `json:"operator"`
}

type Operator string

const (
	OperatorIs    = Operator("is")
	OperatorIsNot = Operator("is not")

	OperatorExists    = Operator("exists")
	OperatorNotExists = Operator("does not exist")
)
