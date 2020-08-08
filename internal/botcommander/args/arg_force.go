package args

/*
	ARGUMENT: Force
*/

const (
	// ForceDeployName is the key/id for the force deploy argument
	ForceDeployName = "force"
	// ForceDeployDescription is the description of the Force Deploy argument
	ForceDeployDescription = "forces the command to execute"
)

// Force is the Force Deploy argument bool type
type Force bool

// Name is the name of the Force Deploy argument
func (a Force) Name() string {
	return ForceDeployName
}

// Value is the value of the Force Deploy argument
func (a Force) Value() interface{} {
	return bool(a)
}

// Description is the description of the Force Deploy argument
func (a Force) Description() string {
	return ForceDeployDescription
}

// DefaultForceArg is the default Force Deploy argument
func DefaultForceArg() Force {
	return false
}
