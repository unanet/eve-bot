package help

// Examples Evebot Command Examples
type Examples []string

// String converts example slice to a string
func (ebce Examples) String() string {
	var msg string
	for _, s := range ebce {
		if len(msg) > 0 {
			msg = msg + "\n" + "@evebot " + s
		} else {
			msg = "@evebot " + s
		}
	}
	return msg
}
