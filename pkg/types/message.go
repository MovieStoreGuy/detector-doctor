package types

// Message defines the generic result returned back from the websocket connection
type Message interface {
	GetType() string
}
