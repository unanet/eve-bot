package botargs

import (
	"strconv"
	"strings"
)

type Arg interface {
	Name() string
	Description() string
}

type Args []Arg

func (a Args) String() string {
	var msg string
	for _, v := range a {
		msg = msg + v.Name() + " - " + v.Description() + "\n"
	}
	return msg
}

func ResolveArgumentKV(argKV []string) Arg {
	switch strings.ToLower(argKV[0]) {
	case "dryrun":
		b, err := strconv.ParseBool(strings.ToLower(argKV[1]))
		if err != nil {
			return nil
		} else {
			return Dryrun(b)
		}
	case "force":
		b, err := strconv.ParseBool(strings.ToLower(argKV[1]))
		if err != nil {
			return nil
		} else {
			return Force(b)
		}
	case "services":
		return NewServicesArg(strings.Split(argKV[1], ","))
	case "databases":
		return NewDatabasesArg(strings.Split(argKV[1], ","))
	default:
		return nil
	}
}
