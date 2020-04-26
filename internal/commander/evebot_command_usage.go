package commander

// Evebot Command Usage
type EvebotCommandUsage []string

func (ebcu EvebotCommandUsage) String() string {
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
