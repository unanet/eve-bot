package commander

// Evebot Command Examples
type UserHelpExamples []string

func (ebce UserHelpExamples) String() string {
	var msg string
	for _, s := range ebce {
		if len(msg) > 0 {
			msg = msg + "\n" + s
		} else {
			msg = s
		}
	}
	return msg
}
