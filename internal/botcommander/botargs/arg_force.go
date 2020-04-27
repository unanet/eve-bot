package botargs

/*
	ARGUMENT: Force
*/

type Force bool

func (a Force) Name() string {
	return "force"
}

func (a Force) Description() string {
	return "forces the command to execute"
}

func DefaultForceArg() Force {
	return false
}
