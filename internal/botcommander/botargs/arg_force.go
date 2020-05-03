package botargs

/*
	ARGUMENT: Force
*/

const (
	ForceDeployName = "force"
)

type Force bool

func (a Force) Name() string {
	return ForceDeployName
}

func (a Force) Value() interface{} {
	return bool(a)
}

func (a Force) Description() string {
	return "forces the command to execute"
}

func DefaultForceArg() Force {
	return false
}
