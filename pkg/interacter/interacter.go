package interacter

type Interacter interface {
	Name() string
	Enabled() bool
	Init()
	Start()
}
