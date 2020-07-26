package params

const (
	FromFeedName = "from_feed"
	ToFeedName   = "to_feed"
)

type Feed struct {
	baseParam
}

func (e Feed) Name() string {
	return e.name
}

func (e Feed) Description() string {
	return e.description
}

func (e Feed) Value() string {
	return e.value
}

func DefaultFromFeed() Feed {
	return Feed{baseParam{
		name:        FromFeedName,
		description: "the name of the source feed",
	}}
}

func DefaultToFeed() Feed {
	return Feed{baseParam{
		name:        ToFeedName,
		description: "the name of the destination feed",
	}}
}
