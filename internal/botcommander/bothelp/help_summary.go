package bothelp

// Evebot Command Summary
type Summary string

func (ebcs Summary) String() string {
	return string(ebcs)
}
