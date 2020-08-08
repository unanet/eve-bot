package params

import (
	"encoding/json"
)

const (
	// MetadataName key/id
	MetadataName = "metadata"
)

// Metadata param data struct
type Metadata struct {
	baseParam
}

// Name satisfies the param interface and returns the Metadata Name
func (e Metadata) Name() string {
	return e.name
}

// Description satisfies the param interface and returns the Metadata Description
func (e Metadata) Description() string {
	return e.description
}

// Value satisfies the param interface and returns the Metadata Value
func (e Metadata) Value() string {
	return e.value
}

// DefaultMetadata is the default ToFeed (used for help/init)
func DefaultMetadata() Metadata {
	return Metadata{baseParam{
		name:        MetadataName,
		description: "the metadata for a service",
	}}
}

// MetadataKeys contains a slice of the metadata key
type MetadataKeys []string

// MetadataMap data structure to hold the metadata
type MetadataMap map[string]interface{}

// ToString converts the MetadataMap to a JSON string
func (e MetadataMap) ToString() string {
	if e == nil || len(e) == 0 {
		return "no metadata"
	}
	jsonB, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		return "invalid json metadata"
	}
	return string(jsonB)
}
