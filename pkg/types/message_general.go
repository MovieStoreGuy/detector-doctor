package types

import "encoding/json"

const ExpectedGeneralMessageType = "message"

type MessageGeneral struct {
	Kind               string                 `json:"type"`
	Channel            string                 `json:"channel"`
	LogicalTimeStampMs int64                  `json:"logicalTimestampMs"`
	Message            map[string]interface{} `json:"message"`
}

func (gen *MessageGeneral) GetType() string {
	return gen.Kind
}

func ReadGeneralMessage(data []byte) *MessageGeneral {
	var gen MessageGeneral
	if err := json.Unmarshal(data, &gen); err == nil && ExpectedGeneralMessageType == gen.Kind {
		return &gen
	}
	return nil
}
