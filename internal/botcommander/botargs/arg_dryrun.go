package botargs

/*
	ARGUMENT: Dryrun
*/

const (
	DryrunName = "dryrun"
)

type Dryrun bool

func (a Dryrun) Name() string {
	return DryrunName
}

func (a Dryrun) Value() interface{} {
	return bool(a)
}

func (a Dryrun) Description() string {
	return "generates a plan but doesn't actually change any state"
}

func DefaultDryrunArg() Dryrun {
	return false
}
