package ports

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/interactivehub/engine/app"
	"github.com/interactivehub/engine/app/command"
	"github.com/interactivehub/engine/domain/user"
	"github.com/pkg/errors"
)

type UsersGrpcServer struct {
	user.UnimplementedUsersServiceServer
	app app.Application
}

func NewUsersGrpcServer(app app.Application) *UsersGrpcServer {
	return &UsersGrpcServer{app: app}
}

func (s *UsersGrpcServer) NewUser(ctx context.Context, req *user.NewUserRequest) (*empty.Empty, error) {
	cmd := command.NewUser{
		ID:       req.GetUserId(),
		UniqueID: req.GetUniqueId(),
		Nickname: req.GetNickname(),
	}

	err := s.app.Commands.NewUser.Handle(ctx, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	return &empty.Empty{}, nil
}
