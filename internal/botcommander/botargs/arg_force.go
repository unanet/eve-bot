package botargs

/*
	ARGUMENT: Force
*/

type ArgForce bool

func (ebaf ArgForce) Name() string {
	return "force"
}

func (ebaf ArgForce) Description() string {
	return "forces the command to execute"
}

func NewForceArg() ArgForce {
	return false
}
