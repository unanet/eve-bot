package botargs

/*
	ARGUMENT: Databases
*/

type Databases []string

func (a Databases) Name() string {
	return "databases"
}

func (a Databases) Description() string {
	return "comma separated list of databases"
}

func DefaultDatabasesArg() Databases {
	return Databases{}
}

func NewDatabasesArg(input []string) Databases {
	var dbs = Databases{}

	for _, v := range input {
		dbs = append(dbs, v)
	}

	return dbs
}
