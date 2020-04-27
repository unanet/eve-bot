package commander

// Evebot Command Usage
type HelpUsage []string

func (ebcu HelpUsage) String() string {
	var msg string
	for _, s := range ebcu {
		if len(msg) > 0 {
			msg = msg + "\n" + s
		} else {
			msg = s
		}
	}
	return msg
}
