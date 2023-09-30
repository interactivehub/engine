package ports

import (
	"context"

	"github.com/interactivehub/engine/app"
	"github.com/interactivehub/engine/app/handlers"
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

func (s *UsersGrpcServer) NewUser(ctx context.Context, req *user.NewUserRequest) (*user.NewUserResponse, error) {
	cmd := handlers.NewUser{
		ID:       req.GetUserId(),
		UniqueID: req.GetUniqueId(),
		Nickname: req.GetNickname(),
	}

	res, err := s.app.Handlers.NewUser.Handle(ctx, cmd)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	return &user.NewUserResponse{UserId: res.ID}, nil
}
