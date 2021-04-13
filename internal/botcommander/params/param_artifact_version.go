package params

const (
	// ArtifactVersionName is the key/id for the ArtifactVersion Param
	ArtifactVersionName = "version"
)

// ArtifactVersion data structure
type ArtifactVersion struct {
	baseParam
}

// Name satisfies the param interface and returns the Artifact Version Name
func (e ArtifactVersion) Name() string {
	return e.name
}

// Description satisfies the param interface and returns the Artifact Version Description
func (e ArtifactVersion) Description() string {
	return e.description
}

// Value satisfies the param interface and returns the Artifact Version Value
func (e ArtifactVersion) Value() string {
	return e.value
}
