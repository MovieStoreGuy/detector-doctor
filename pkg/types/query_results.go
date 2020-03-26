package types

import "encoding/json"

// QueryResults defines the abstracted bulk return from SignalFx when asking for items
// in bulk, the results field is left as json.RawMessages to allow for partial proccessing
type QueryResults struct {
	Count   int32             `json:"count"`
	Results []json.RawMessage `json:"results"`
}
