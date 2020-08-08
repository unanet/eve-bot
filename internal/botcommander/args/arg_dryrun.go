package args

/*
	ARGUMENT: Dryrun
*/

const (
	// DryrunName is used a key/id for the dryrun option
	DryrunName = "dryrun"
	// DryrunDescription is the description for the dryrun
	DryrunDescription = "generates a plan but doesn't actually change any state"
)

// Dryrun is the dryun bool type
type Dryrun bool

// Name is the name of the dryrun argument
func (a Dryrun) Name() string {
	return DryrunName
}

// Value is the value of the dryrun argument
func (a Dryrun) Value() interface{} {
	return bool(a)
}

// Description is the description of the dryrun
func (a Dryrun) Description() string {
	return DryrunDescription
}

// DefaultDryrunArg is the default dryrun argument
func DefaultDryrunArg() Dryrun {
	return false
}
