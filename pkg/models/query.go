package models

type Query []Filter

type Filter []Condition

type Condition struct {
	Type   ConditionType          `json:"type"`
	Params map[string]interface{} `json:"params"`
}

type ConditionType string

const (
	DevicePropertyCondition       = ConditionType("DevicePropertyCondition")
	LabelValueCondition           = ConditionType("LabelValueCondition")
	LabelExistenceCondition       = ConditionType("LabelExistenceCondition")
	ApplicationReleaseCondition   = ConditionType("ApplicationReleaseCondition")
	ApplicationExistenceCondition = ConditionType("ApplicationExistenceCondition")
	ServiceStateCondition         = ConditionType("ServiceStateCondition")
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

const (
	AnyApplicationRelease = string("any")
)

type ApplicationReleaseConditionParams struct {
	ApplicationID string   `json:"applicationId"`
	Operator      Operator `json:"operator"`
	Release       string   `json:"release"`
}

type ApplicationExistenceConditionParams struct {
	ApplicationID string   `json:"applicationId"`
	Operator      Operator `json:"operator"`
}

type ServiceStateConditionParams struct {
	ApplicationID string       `json:"applicationId"`
	Service       string       `json:"service"`
	Operator      Operator     `json:"operator"`
	ServiceState  ServiceState `json:"serviceState"`
}

type Operator string

const (
	OperatorIs    = Operator("is")
	OperatorIsNot = Operator("is not")

	OperatorExists    = Operator("exists")
	OperatorNotExists = Operator("does not exist")
)
