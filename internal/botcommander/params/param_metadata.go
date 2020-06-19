package params

import "fmt"

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
	msg := ""
	for i, v := range e {
		msg += fmt.Sprintf("%s=%s", i, v)
	}
	return msg
}
