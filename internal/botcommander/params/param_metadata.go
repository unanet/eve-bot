package params

import (
	"encoding/json"
)

const (
	MetadataName = "metadata"
)

type Metadata struct {
	baseParam
}

func (e Metadata) Name() string {
	return e.name
}

func (e Metadata) Description() string {
	return e.description
}

func (e Metadata) Value() string {
	return e.value
}

func DefaultMetadata() Namespace {
	return Namespace{baseParam{
		name:        MetadataName,
		description: "the metadata for a service",
	}}
}

type MetadataMap map[string]interface{}

func (e MetadataMap) ToString() string {
	jsonB, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		return "invalid json metadata"
	}
	return "```" + string(jsonB) + "```"
}
