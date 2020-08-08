package help

// Summary Evebot Command Summary
type Summary string

// String converts the Summary help to a string
func (ebcs Summary) String() string {
	return string(ebcs)
}
