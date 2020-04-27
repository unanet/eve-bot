package botargs

/*
	ARGUMENT: Dryrun
*/

type Dryrun bool

func (a Dryrun) Name() string {
	return "dryrun"
}

func (a Dryrun) Description() string {
	return "generates a plan but doesn't actually change any state"
}

func DefaultDryrunArg() Dryrun {
	return false
}
