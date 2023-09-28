package adapters

import (
	"context"
	"errors"
	"log"

	"github.com/interactivehub/engine/domain/user"
)

type UsersGrpcServer struct {
	user.UnimplementedUsersServiceServer
	usersRepo user.Repository
}

func NewUsersGrpcServer(usersRepo user.Repository) *UsersGrpcServer {
	if usersRepo == nil {
		panic("missing usersRepo")
	}

	return &UsersGrpcServer{usersRepo: usersRepo}
}

func (u *UsersGrpcServer) NewUser(ctx context.Context, req *user.NewUserRequest) (*user.NewUserResponse, error) {
	log.Println("Received req: ", req.GetNickname())

	userExists, err := u.usersRepo.UserWithIdExists(ctx, req.UserId)
	if err != nil {
		return &user.NewUserResponse{}, err
	}

	if userExists {
		return &user.NewUserResponse{}, errors.New("failed to create user. user already exists")
	}

	usr, err := user.NewUser(
		req.GetUserId(),
		req.GetUniqueId(),
		req.GetNickname(),
		0,
	)
	if err != nil {
		return &user.NewUserResponse{}, err
	}

	err = u.usersRepo.CreateUser(ctx, *usr)
	if err != nil {
		return &user.NewUserResponse{}, err
	}

	return &user.NewUserResponse{
		UserId: usr.ID(),
	}, nil
}
