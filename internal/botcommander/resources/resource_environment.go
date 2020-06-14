package resources

const (
	EnvironmentName = "environments"
)

type Environment struct {
	baseResource
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
	return Environment{baseResource{
		name:        EnvironmentName,
		description: "the environment resource (i.e. una-dev,una-int,una-qa)",
	}}
}
