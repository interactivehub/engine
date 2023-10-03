package ports

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/interactivehub/engine/app"
	"github.com/interactivehub/engine/common/ws"
)

type WSListener interface {
	ListenEvents()
}

type wsListener struct {
	client *websocket.Conn
	app    app.Application
}

func NewWSListener(client *websocket.Conn, app app.Application) *wsListener {
	return &wsListener{client, app}
}

func (l *wsListener) ListenEvents() {
	for {
		var event ws.Event

		err := l.client.ReadJSON(&event)
		if err != nil {
			log.Println(err)
			break
		}

		log.Println(event)
	}
}
