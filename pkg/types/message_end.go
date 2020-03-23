package types

import "encoding/json"

const (
	ExpectControlType = "control-message"

	ControlMessageEndEvent         = "END_OF_CHANNEL"
	ControlMessageStartStreamEvent = "STREAM_START"
)

type ControlMessage struct {
	Event       string `json:"event"`
	Kind        string `json:"type"`
	Channel     string `json:"channel"`
	TimestampMS int64  `json:"timestamp"`
	Progress    int    `json:"progress,omitempty`
}

func ReadControlMessage(data []byte) *ControlMessage {
	var control ControlMessage
	if err := json.Unmarshal(data, &control); err == nil && control.Kind == ExpectControlType {
		return &control
	}
	return nil
}

func IsEndofMessages(msg *ControlMessage) bool {
	return msg != nil && ControlMessageEndEvent == msg.Event
}
