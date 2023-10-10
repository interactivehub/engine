package ports

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/interactivehub/engine/app"
	"github.com/interactivehub/engine/app/command"
	"github.com/interactivehub/engine/domain/games/wheel"
	"github.com/pkg/errors"
)

type WheelGrpcServer struct {
	wheel.UnimplementedWheelServiceServer
	app app.Application
}

func NewWheelGrpcServer(app app.Application) *WheelGrpcServer {
	return &WheelGrpcServer{app: app}
}

func (s *WheelGrpcServer) JoinWheelRound(ctx context.Context, req *wheel.JoinWheelRoundRequest) (*empty.Empty, error) {
	pick, err := wheel.ParseWheelItemColor(req.GetPick().String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse wheel item color pick")
	}

	cmd := command.JoinWheelRound{
		UserID: req.GetUserId(),
		Bet:    req.GetBet(), // TODO: Check this bitch ass bet that isnt parsed correctly
		Pick:   pick,
	}

	err = s.app.Commands.JoinWheelRound.Handle(ctx, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to join wheel round")
	}

	return &empty.Empty{}, nil
}
