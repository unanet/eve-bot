package commander

// EvebotCommand interface
// each evebot command needs to implement this interface
type EvebotCommand interface {
	Name() string
	Help() *EvebotCommandHelp
	Initialize(input []string) EvebotCommand
	IsHelpRequest() bool
	AdditionalArgs() (EvebotCommandArgs, error)
}
