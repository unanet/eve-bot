package params

const (
	ArtifactVersionName = "version"
)

type ArtifactVersion struct {
	baseParam
}

func (e ArtifactVersion) Name() string {
	return e.name
}

func (e ArtifactVersion) Description() string {
	return e.description
}

func (e ArtifactVersion) Value() string {
	return e.value
}

func DefaultArtifactVersion() ArtifactVersion {
	return ArtifactVersion{baseParam{
		name:        ArtifactVersionName,
		description: "the version of the artifact",
	}}
}
