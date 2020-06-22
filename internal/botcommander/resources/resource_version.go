package resources

const (
	VersionName = "version"
)

type Version struct {
	baseResource
}

func (e Version) Name() string {
	return e.name
}

func (e Version) Description() string {
	return e.description
}

func (e Version) Value() string {
	return e.value
}

func DefaultVersion() Version {
	return Version{baseResource{
		name:        VersionName,
		description: "the version resource",
	}}
}
