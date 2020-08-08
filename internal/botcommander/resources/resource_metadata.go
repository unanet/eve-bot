package resources

const (
	// MetadataName resource key/id
	MetadataName = "metadata"
)

// Metadata resource data structure
type Metadata struct {
	baseResource
}

// Name satisfies the resource interface and returns the Metadata Name
func (e Metadata) Name() string {
	return e.name
}

// Description satisfies the resource interface and returns the Metadata Description
func (e Metadata) Description() string {
	return e.description
}

// Value satisfies the resource interface and returns the Metadata Value
func (e Metadata) Value() string {
	return e.value
}

// DefaultMetadata returns the Default Metadata
func DefaultMetadata() Metadata {
	return Metadata{baseResource{
		name:        MetadataName,
		description: "the metadata resource (ex: the config for a service)",
	}}
}
