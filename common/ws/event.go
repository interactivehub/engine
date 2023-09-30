package ws

import (
	"encoding/json"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func NewEvent(t string, p interface{}) Event {
	if t == "" {
		panic("missing event type")
	}

	jsonBytes, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}

	rawMessage := json.RawMessage(jsonBytes)

	return Event{
		Type:    t,
		Payload: rawMessage,
	}
}
