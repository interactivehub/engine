package ports

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/interactivehub/engine/app"
	"github.com/interactivehub/engine/app/command"
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

		handleEvent(context.Background(), l.app, event)
	}
}

const (
	RequestWheelRoundEventType = "requestWheelRound"
)

type RequestWheelRoundEventPayload struct {
	ClientSeed string `json:"clientSeed"`
}

func handleEvent(ctx context.Context, app app.Application, event ws.Event) {
	switch event.Type {
	case RequestWheelRoundEventType:
		var payload RequestWheelRoundEventPayload
		json.Unmarshal(event.Payload, &payload)

		app.Commands.StartWheelRound.Handle(ctx, command.StartWheelRound{ClientSeed: []byte(payload.ClientSeed)})
	}
}
