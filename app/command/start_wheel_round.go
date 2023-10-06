package command

import (
	"context"

	"github.com/interactivehub/engine/adapters"
	"github.com/interactivehub/engine/common/decorator"
	"github.com/patrickmn/go-cache"
)

type StartWheelRound struct{}

type StartWheelRoundHandler decorator.CommandHandler[StartWheelRound]

type startWheelRoundHandler struct {
	memCache *cache.Cache
	wsWriter adapters.WSWriter
}

func NewStartWheelRoundHandler(
	memCache *cache.Cache,
	wsWriter adapters.WSWriter,
) StartWheelRoundHandler {
	if memCache == nil {
		panic("missing memCache")
	}

	if wsWriter == nil {
		panic("missing wsWriter")
	}

	return startWheelRoundHandler{memCache, wsWriter}
}

func (h startWheelRoundHandler) Handle(ctx context.Context, cmd StartWheelRound) error {
	return nil
}
