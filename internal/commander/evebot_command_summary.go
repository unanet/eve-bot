package commander

// Evebot Command Summary
type EvebotCommandSummary string

func (ebcs EvebotCommandSummary) String() string {
	return string(ebcs)
}
