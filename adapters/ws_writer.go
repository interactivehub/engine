package adapters

import (
	"github.com/gorilla/websocket"
	"github.com/interactivehub/engine/common/ws"
	"github.com/pkg/errors"
)

type WSWriter interface {
	WriteEvent(t string, p interface{}) error
}

type wsWriter struct {
	client *websocket.Conn
}

func NewWSWriter(client *websocket.Conn) *wsWriter {
	if client == nil {
		panic("missing client")
	}

	return &wsWriter{client}
}

func (w *wsWriter) WriteEvent(t string, p interface{}) error {
	if w.client == nil {
		return errors.New("missing ws client")
	}

	event := ws.NewEvent(t, p)

	return w.client.WriteJSON(event)
}
