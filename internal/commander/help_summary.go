package commander

// Evebot Command Summary
type HelpSummary string

func (ebcs HelpSummary) String() string {
	return string(ebcs)
}
