package help

// Usage Evebot Command Usage
type Usage []string

// String converts the User string slice to a string
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
