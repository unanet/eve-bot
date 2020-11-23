package resources

const (
	// MetadataName resource key/id
	JobName = "job"
)

// Metadata resource data structure
type Job struct {
	baseResource
}

// Name satisfies the resource interface and returns the Metadata Name
func (e Job) Name() string {
	return e.name
}

// Description satisfies the resource interface and returns the Metadata Description
func (e Job) Description() string {
	return e.description
}

// Value satisfies the resource interface and returns the Metadata Value
func (e Job) Value() string {
	return e.value
}

// DefaultMetadata returns the Default Metadata
func DefaultJob() Job {
	return Job{baseResource{
		name:        JobName,
		description: "the job to run",
	}}
}
