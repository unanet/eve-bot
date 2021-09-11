package commands

type authCmd struct {
	baseCommand
}

const (
	// AuthCmdName used as key/id for the auth command
	AuthCmdName = "auth"
)

// NewAuthCommand creates a New AuthCmd that implements the EvebotCommand interface
func NewAuthCommand(cmdFields []string, channel, user string) EvebotCommand {
	return authCmd{baseCommand{
		input: cmdFields,
		info: ChatInfo{
			User:          user,
			Channel:       channel,
			CommandName:   AuthCmdName,
			IsHelpRequest: isHelpCmd(cmdFields, AuthCmdName),
			IsAuthCmd:     true,
		},
		opts:   make(CommandOptions),
		bounds: InputLengthBounds{Min: -1, Max: -1},
	}}
}

// AckMsg satisfies the EveBotCommand Interface and returns the acknowledgement message
// ...and whether or not we should continue
func (cmd authCmd) AckMsg() (string, bool) {
	return "Please Check your Private DM from `evebot` for an auth link", true
}

// Options satisfies the EveBotCommand Interface and returns the dynamic options
func (cmd authCmd) Options() CommandOptions {
	return cmd.opts
}

// Info satisfies the EveBotCommand Interface and returns the Chat Info
func (cmd authCmd) Info() ChatInfo {
	return cmd.info
}
