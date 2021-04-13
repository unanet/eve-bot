package resources

const (
	// EnvironmentName resource key/id
	EnvironmentName = "environments"
)

// Environment resource data structure
type Environment struct {
	baseResource
}

// Name satisfies the resource interface and returns the Environment Name
func (e Environment) Name() string {
	return e.name
}

// Description satisfies the resource interface and returns the Environment Description
func (e Environment) Description() string {
	return e.description
}

// Value satisfies the resource interface and returns the Environment Value
func (e Environment) Value() string {
	return e.value
}
