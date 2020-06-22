package params

const (
	VersionName = "version"
)

type Version struct {
	baseParam
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
	return Version{baseParam{
		name:        VersionName,
		description: "the version param (ex: 20.1, 20.2, 0.0)",
	}}
}
