package engine

type Instance struct {
	ID      string
	Labels  map[string]string
	Running bool
}
