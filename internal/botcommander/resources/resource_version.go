package resources

const (
	// VersionName resource key/id
	VersionName = "version"
)

// Version resource data structure
type Version struct {
	baseResource
}

// Name satisfies the resource interface and returns the Version Name
func (e Version) Name() string {
	return e.name
}

// Description satisfies the resource interface and returns the Version Description
func (e Version) Description() string {
	return e.description
}

// Value satisfies the resource interface and returns the Version Value
func (e Version) Value() string {
	return e.value
}
