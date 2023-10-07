package decorator

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type commandLoggingDecorator[C any] struct {
	base   CommandHandler[C]
	logger *zap.Logger
}

func (d commandLoggingDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	handlerType := generateActionName(cmd)

	logger := d.logger.With(
		zap.String("command", handlerType),
		zap.String("command_body", fmt.Sprintf("%#v", cmd)),
	)

	logger.Debug("Executing command")
	defer func() {
		if err == nil {
			logger.Info("Command executed successfully")
		} else {
			logger.Error(errors.Wrap(err, "Failed to execute command").Error())
		}
	}()

	return d.base.Handle(ctx, cmd)
}
