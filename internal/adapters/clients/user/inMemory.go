package user

import (
	"context"

	user "github.com/sopial42/cleanic/internal/domains/user"
	userSVC "github.com/sopial42/cleanic/internal/services/user"
	authSVC "github.com/sopial42/cleanic/internal/services/auth"
)

type inMemory struct {
	userSvc userSVC.Service
}

func NewInMemoryUserClient(userSvc userSVC.Service) authSVC.UserClient {
	return &inMemory{
		userSvc: userSvc,
	}
}

func (m *inMemory) Create(ctx context.Context, newUser user.User) (user.User, error) {
	return m.userSvc.Create(ctx, newUser)
}

func (m *inMemory) GetUserByEmail(ctx context.Context, email user.Email) (user.User, error) {
	return m.userSvc.GetUserByEmail(ctx, email)
}
