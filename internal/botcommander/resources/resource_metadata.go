package resources

const (
	MetadataName = "metadata"
)

type Metadata struct {
	baseResource
}

func (e Metadata) Name() string {
	return e.name
}

func (e Metadata) Description() string {
	return e.description
}

func (e Metadata) Value() string {
	return e.value
}

func DefaultMetadata() Metadata {
	return Metadata{baseResource{
		name:        MetadataName,
		description: "the metadata resource (ex: the config for a service)",
	}}
}
