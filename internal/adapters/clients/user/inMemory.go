package user

import (
	"context"

	user "github.com/sopial42/cleanic/internal/domains/user"
	authSVC "github.com/sopial42/cleanic/internal/services/auth"
	userSVC "github.com/sopial42/cleanic/internal/services/user"
)

type inMemory struct {
	userSVC userSVC.Service
}

func NewInMemoryUserClient(userSVC userSVC.Service) authSVC.UserClient {
	return &inMemory{
		userSVC: userSVC,
	}
}

func (m *inMemory) Create(ctx context.Context, newUser user.User) (user.User, error) {
	return m.userSVC.Create(ctx, newUser)
}

func (m *inMemory) GetUserByEmail(ctx context.Context, email user.Email) (user.User, error) {
	return m.userSVC.GetUserByEmail(ctx, email)
}

func (m *inMemory) GetUserByID(ctx context.Context, userID user.ID) (user.User, error) {
	return m.userSVC.GetUserByID(ctx, userID)
}
