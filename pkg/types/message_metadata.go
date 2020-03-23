package types

import "encoding/json"

const ExpectedMetadataType = "metadata"

type MessageMetadata struct {
	TimeSeriesID string `json:"tsId"`
	Kind string `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Channel string `json:"channel"`
}

func (meta *MessageMetadata) GetType() string {
	return meta.Kind
}

func ReadMetadataMessage(data []byte) *MessageMetadata {
	var metadata MessageMetadata
	if err := json.Unmarshal(data, &metadata); err == nil && metadata.Kind == ExpectedMetadataType {
		return &metadata
	}
	return nil
}