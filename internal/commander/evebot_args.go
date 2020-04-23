package commander

type EvebotArg interface {
	Name() string
	//Type() reflect.Type
}

type EvebotArgDryrun bool
type EvebotArgForce bool
type EvebotArgServices []string
type EvebotArgDatabases []string
type EvebotArgs []EvebotArg

func (ebad EvebotArgDryrun) Name() string {
	return "dryrun"
}

//func (ebad EvebotArgDryrun) Type() reflect.Type {
//	return reflect.TypeOf(ebad)
//}

func (ebaf EvebotArgForce) Name() string {
	return "force"
}

func (ebas EvebotArgServices) Name() string {
	return "services"
}

func (ebad EvebotArgDatabases) Name() string {
	return "databases"
}

func NewAdditionArg(argT interface{}, val interface{}) EvebotArg {
	switch argT.(type) {
	case EvebotArgDryrun:
		if b, ok := val.(bool); ok {
			return EvebotArgDryrun(b)
		}
	case EvebotArgForce:
		if b, ok := val.(bool); ok {
			return EvebotArgForce(b)
		}
	case EvebotArgServices:
		if b, ok := val.([]string); ok {
			return EvebotArgServices(b)
		}
	case EvebotArgDatabases:
		if b, ok := val.([]string); ok {
			return EvebotArgDatabases(b)
		}
	}
	return nil
}
