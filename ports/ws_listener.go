package ports

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/interactivehub/engine/common/ws"
)

type WSListener interface {
	ListenEvents()
}

type wsListener struct {
	client *websocket.Conn
}

func NewWSListener(client *websocket.Conn) *wsListener {
	return &wsListener{client}
}

func (l *wsListener) ListenEvents() {
	for {
		var event ws.Event

		err := l.client.ReadJSON(&event)
		if err != nil {
			log.Println(err)
		}

		log.Println(event)
	}
}
