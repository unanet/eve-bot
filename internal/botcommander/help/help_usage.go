package help

// Evebot Command Usage
type Usage []string

func (ebcu Usage) String() string {
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
