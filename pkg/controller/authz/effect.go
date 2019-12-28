package authz

type Effect string

const (
	EffectAllow = Effect("allow")
	EffectDeny  = Effect("deny")
)
