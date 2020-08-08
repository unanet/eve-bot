package params

const (
	// VersionName param key/id
	VersionName = "version"
)

// Version param data struct
type Version struct {
	baseParam
}

// Name satisfies the param interface and returns the Version Name
func (e Version) Name() string {
	return e.name
}

// Description satisfies the param interface and returns the Version Description
func (e Version) Description() string {
	return e.description
}

// Value satisfies the param interface and returns the Version Value
func (e Version) Value() string {
	return e.value
}

// DefaultVersion is the default Version (used for help/init)
func DefaultVersion() Version {
	return Version{baseParam{
		name:        VersionName,
		description: "the version param (ex: 20.1, 20.2, 0.0)",
	}}
}
