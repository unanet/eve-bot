package commander

/*
	ARGUMENT: Force
*/

type EvebotArgForce bool

func (ebaf EvebotArgForce) Name() string {
	return "force"
}

func (ebaf EvebotArgForce) Description() string {
	return "forces the command to execute"
}

func NewForceArg() EvebotArgForce {
	return false
}
