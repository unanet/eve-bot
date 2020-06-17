package eveapimodels

type ArtifactDefinition struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	RequestedVersion string `json:"requested_version,omitempty"`
	AvailableVersion string `json:"available_version"`
	ArtifactoryFeed  string `json:"artifactory_feed"`
	ArtifactoryPath  string `json:"artifactory_path"`
	FunctionPointer  string `json:"function_pointer"`
	FeedType         string `json:"feed_type"`
	Matched          bool   `json:"-"`
}

type ArtifactDefinitions []*ArtifactDefinition
