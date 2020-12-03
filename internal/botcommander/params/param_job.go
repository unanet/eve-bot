package params

const (
	// JobName param key/id
	JobName = "job"
)

// Job param data struct
type Job struct {
	baseParam
}

// Name satisfies the param interface and returns the Namespace Name
func (e Job) Name() string {
	return e.name
}

// Description satisfies the param interface and returns the Namespace Description
func (e Job) Description() string {
	return e.description
}

// Value satisfies the param interface and returns the Namespace Value
func (e Job) Value() string {

	return e.value
}

// DefaultNamespace is the default Namespace (used for help/init)
func DefaultJob() Job {
	return Job{baseParam{
		name:        JobName,
		description: "the job to run",
	}}
}
