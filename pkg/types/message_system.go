package types

import "encoding/json"

const (
	ExeceptedKeepAliveEvent = "KEEP_ALIVE"
)

type MessageKeepAlive struct {
	Event       string `json:"event"`
	TimestampMS int64  `json:"timestampMs"`
}

func ReadKeepAliveMessage(data []byte) *MessageKeepAlive {
	var keepalive MessageKeepAlive
	if err := json.Unmarshal(data, &keepalive); err == nil && keepalive.Event == ExeceptedKeepAliveEvent {
		return &keepalive
	}
	return nil
}
