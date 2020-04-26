package commander

// Evebot Command Examples
type EvebotCommandExamples []string

func (ebce EvebotCommandExamples) String() string {
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
