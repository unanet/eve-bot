package commander

// EvebotCommand interface
// each evebot command needs to implement this interface
type EvebotCommand interface {
	Name() string
	Help() *Help
	Initialize(input []string) EvebotCommand
	IsHelpRequest() bool
	AdditionalArgs() (Args, error)
	AsyncRequired() bool
}
