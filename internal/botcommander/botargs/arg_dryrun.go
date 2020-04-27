package botargs

/*
	ARGUMENT: Dryrun
*/

type ArgDryrun bool

func (ebad ArgDryrun) Name() string {
	return "dryrun"
}

func (ebad ArgDryrun) Description() string {
	return "generates a plan but doesn't actually change any state"
}

func NewDryrunArg() ArgDryrun {
	return false
}
