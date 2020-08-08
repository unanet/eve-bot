package params

const (
	// ArtifactName is the key/id for the Artifact Param
	ArtifactName = "artifact"
)

// Artifact is the data structure for the Artifact Param
type Artifact struct {
	baseParam
}

// Name satisfies the param interface and returns the artifact Name
func (e Artifact) Name() string {
	return e.name
}

// Description satisfies the param interface and returns the artifact Description
func (e Artifact) Description() string {
	return e.description
}

// Value satisfies the param interface and returns the artifact Value
func (e Artifact) Value() string {
	return e.value
}

// DefaultArtifact is the default Artifact (used for help/init)
func DefaultArtifact() Artifact {
	return Artifact{baseParam{
		name:        ArtifactName,
		description: "the name of the artifact",
	}}
}
