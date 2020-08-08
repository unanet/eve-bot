package params

const (
	// EnvironmentName is the key/id for the Environment Param
	EnvironmentName = "environment"
)

// Environment param data struct
type Environment struct {
	baseParam
}

// Name satisfies the param interface and returns the Environment Name
func (e Environment) Name() string {
	return e.name
}

// Description satisfies the param interface and returns the Environment Description
func (e Environment) Description() string {
	return e.description
}

// Value satisfies the param interface and returns the Environment Value
func (e Environment) Value() string {
	return e.value
}

// DefaultEnvironment is the default Environment (used for help/init)
func DefaultEnvironment() Environment {
	return Environment{baseParam{
		name:        EnvironmentName,
		description: "the environment where the resource lives (i.e. dev,int,qa,stage,perf,prod)",
	}}
}
