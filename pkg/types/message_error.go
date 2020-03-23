package types

import (
	"encoding/json"
	"fmt"
)

type WebsocketError struct {
	ErrorCode int                    `json:"error"`
	ErrorType string                 `json:"errorType"`
	Message   string                 `json:"message"`
	Channel   string                 `json:"channel"`
	Context   map[string]interface{} `json:"context"`
	Kind      string                 `json:"type"`
}

func (web *WebsocketError) Error() string {
	return fmt.Sprintf("code:%v type:%v msg:%s", web.ErrorCode, web.ErrorType, web.Message)
}

func IsWebsocketError(err error) bool {
	_, cast := err.(*WebsocketError)
	return cast
}

func ReadWebsocketError(data []byte) *WebsocketError {
	var web WebsocketError
	if err := json.Unmarshal(data, &web); err == nil && web.Kind == "error" {
		return &web
	}
	return nil
}
