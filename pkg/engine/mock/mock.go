package mock

type Engine struct {
	instances engine.Instances
}

type container struct {
	id      string
	running bool
}
