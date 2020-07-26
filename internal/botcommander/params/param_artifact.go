package params

const (
	ArtifactName = "artifact"
)

type Artifact struct {
	baseParam
}

func (e Artifact) Name() string {
	return e.name
}

func (e Artifact) Description() string {
	return e.description
}

func (e Artifact) Value() string {
	return e.value
}

func DefaultArtifact() Artifact {
	return Artifact{baseParam{
		name:        ArtifactName,
		description: "the name of the artifact",
	}}
}
