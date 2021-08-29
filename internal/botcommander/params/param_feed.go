package params

const (
	// FromFeedName is the key/id for the Environment Param
	FromFeedName = "from_feed"
	// ToFeedName is the key/id for the Environment Param
	ToFeedName = "to_feed"
)

// Feed data structure for the Feed Param
type Feed struct {
	baseParam
}

// Name satisfies the param interface and returns the Feed Name
func (e Feed) Name() string {
	return e.name
}

// Description satisfies the param interface and returns the Feed Description
func (e Feed) Description() string {
	return e.description
}

// Value satisfies the param interface and returns the Feed Value
func (e Feed) Value() string {
	return e.value
}
