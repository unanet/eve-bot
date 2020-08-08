package args

import (
	"strconv"
	"strings"
)

// Arg is the interface for command arguments
type Arg interface {
	Name() string
	Description() string
	Value() interface{}
}

// Args are a slice of argument interfaces
type Args []Arg

// String iterates the arguments and concats the Name/Description
func (a Args) String() string {
	var msg string
	for _, v := range a {
		msg = msg + v.Name() + " - " + v.Description() + "\n"
	}
	return msg
}

// ResolveArgumentKV resolves the Key Values
func ResolveArgumentKV(argKV []string) Arg {
	switch strings.ToLower(argKV[0]) {
	case DryrunName:
		b, err := strconv.ParseBool(strings.ToLower(argKV[1]))
		if err != nil {
			return nil
		}
		return Dryrun(b)
	case ForceDeployName:
		b, err := strconv.ParseBool(strings.ToLower(argKV[1]))
		if err != nil {
			return nil
		}
		return Force(b)
	case ServicesName:
		return NewServicesArg(strings.Split(argKV[1], ","))
	case DatabasesName:
		return NewDatabasesArg(strings.Split(argKV[1], ","))
	default:
		return nil
	}
}
