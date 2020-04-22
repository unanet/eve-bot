package commander

// EvebotCommand interface
// each evebot command needs to implement this interface
type EvebotCommand interface {
	Execute() error
}

// EvebotResolver resolves the input commands
type EvebotResolver interface {
	Resolve(input []string) (EvebotCommand, error)
	Help() string
}

// type baseCmd struct {
// 	Name     string
// 	Template string
// 	Examples string
// }
