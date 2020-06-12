package params

const (
	EnvironmentName = "environment"
)

type Environment struct {
	baseParam
}

func (e Environment) Name() string {
	return e.name
}

func (e Environment) Description() string {
	return e.description
}

func (e Environment) Value() string {
	return e.value
}

func DefaultEnvironment() Environment {
	return Environment{baseParam{
		name:        EnvironmentName,
		description: "the environment where the resource lives (i.e. dev,int,qa,stage,perf,prod)",
	}}
}

func NewEnvironmentParam(val string) Environment {
	p := DefaultEnvironment()
	p.value = val
	return p
}
