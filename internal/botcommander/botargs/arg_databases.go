package botargs

/*
	ARGUMENT: Databases
*/

type ArgDatabases []string

func (ebad ArgDatabases) Name() string {
	return "databases"
}

func (ebad ArgDatabases) Description() string {
	return "comma separated list of databases"
}

func NewDatabasesArg() ArgDatabases {
	return ArgDatabases{}
}

func NewParsedArgDatabases(input []string) ArgDatabases {
	var dbs = ArgDatabases{}

	for _, v := range input {
		dbs = append(dbs, v)
	}

	return dbs
}
