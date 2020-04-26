package commander

/*
	ARGUMENT: Dryrun
*/

type EvebotArgDryrun bool

func (ebad EvebotArgDryrun) Name() string {
	return "dryrun"
}

func (ebad EvebotArgDryrun) Description() string {
	return "generates a plan but doesn't actually change any state"
}

func NewDryrunArg() EvebotArgDryrun {
	return false
}
