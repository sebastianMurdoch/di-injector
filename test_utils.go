package di_injector

type A struct {
	Dependency0 string
	Dependency1 string `inject:"auto"`
	Dependency2 Runner `inject:"auto"`
}

type B struct {
	Dependency0 interface{}
}

type Runner interface {
	Run() string
}

type RunnerImpl struct {
	Runner
}

type C struct {
	Dependency0 string `inject:"auto"`
}

func (c *C) Run() string {
	return c.Dependency0
}

type D struct {
	Runner
	Dependency0 int `inject:"auto"`
}