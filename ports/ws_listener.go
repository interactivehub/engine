package ports

import (
	"context"
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/interactivehub/engine/app"
	"github.com/interactivehub/engine/app/command"
	"github.com/interactivehub/engine/common/ws"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type WSListener interface {
	ListenEvents()
}

type wsListener struct {
	client *websocket.Conn
	app    app.Application
	logger *zap.Logger
}

func NewWSListener(
	client *websocket.Conn,
	app app.Application,
	logger *zap.Logger,
) *wsListener {
	return &wsListener{client, app, logger}
}

func (l *wsListener) ListenEvents() {
	for {
		var event ws.Event

		err := l.client.ReadJSON(&event)
		if err != nil {
			l.logger.Error(errors.Wrap(err, "failed to parse event payload").Error())
			break
		}

		handleEvent(context.Background(), l.app, event, l.logger)
	}
}

const (
	RequestWheelRoundEventType = "requestWheelRound"
)

type RequestWheelRoundEventPayload struct {
	ClientSeed string `json:"clientSeed"`
}

func handleEvent(ctx context.Context, app app.Application, event ws.Event, logger *zap.Logger) {
	switch event.Type {
	case RequestWheelRoundEventType:
		payload := &RequestWheelRoundEventPayload{}
		json.Unmarshal(event.Payload, payload)

		err := app.Commands.StartWheelRound.Handle(ctx, command.StartWheelRound{ClientSeed: []byte(payload.ClientSeed)})
		if err != nil {
			logger.Error(errors.Wrap(err, "failed to start wheel round").Error())
		}
	}
}
