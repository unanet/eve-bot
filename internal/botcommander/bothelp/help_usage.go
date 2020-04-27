package bothelp

// Evebot Command Usage
type HelpUsage []string

func (ebcu HelpUsage) String() string {
	var msg string
	for _, s := range ebcu {
		if len(msg) > 0 {
			msg = msg + "\n" + "@evebot " + s
		} else {
			msg = "@evebot " + s
		}
	}
	return msg
}
