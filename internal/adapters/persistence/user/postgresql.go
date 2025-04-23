package persistence

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	user "github.com/sopial42/cleanic/internal/domains/user"
	userSVC "github.com/sopial42/cleanic/internal/services/user"
)

type pgPersistence struct {
	clientDB *bun.DB
}

func NewPGClient(client *bun.DB) userSVC.Persistence {
	return &pgPersistence{clientDB: client}
}

func (p *pgPersistence) Insert(ctx context.Context, newUser user.User) (user.User, error) {
	userDAO := userFromDomainToDAO(newUser)

	_, err := p.clientDB.NewInsert().
		Model(&userDAO).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return user.User{}, fmt.Errorf("err: %w", err)
	}

	// ID == 0 means that the insert failed
	if userDAO.ID == 0 {
		return user.User{}, fmt.Errorf("unable to create a new user: %+v", userDAO)
	}

	return userFromDAOToDomain(userDAO), nil
}

func (p *pgPersistence) ListUsers(ctx context.Context) ([]user.User, error) {
	var userDAOs []userDAO

	request := p.clientDB.NewSelect().Model(&userDAOs)
	err := request.Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to list users: %w", err)
	}

	return userFromDAOsToDomains(userDAOs), nil
}

func (p *pgPersistence) GetUserByID(ctx context.Context, id user.ID) (user.User, error) {
	var userDAO userDAO

	err := p.clientDB.NewSelect().
		Model(&userDAO).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return user.User{}, fmt.Errorf("unable to get user by ID: %w", err)
	}

	if userDAO.ID == 0 {
		return user.User{}, fmt.Errorf("unable to get user by ID: %v", id)
	}

	return userFromDAOToDomain(userDAO), nil
}

func (p *pgPersistence) GetUserByEmail(ctx context.Context, email user.Email) (user.User, error) {
	var userDAO userDAO

	err := p.clientDB.NewSelect().
		Model(&userDAO).
		Where("email = ?", email).
		Scan(ctx)
	if err != nil {
		return user.User{}, fmt.Errorf("unable to get user by email: %w", err)
	}

	if userDAO.ID == 0 {
		return user.User{}, fmt.Errorf("unable to get user by email: %v", email)
	}

	return userFromDAOToDomain(userDAO), nil
}

// UpdateUser perform basic updates on a role but never update roles or ID for safety purpose
func (p *pgPersistence) UpdateUser(ctx context.Context, updatedUser user.User) (user.User, error) {
	userDAO := userFromDomainToDAO(updatedUser)
	userDAO.Roles = []string{}

	_, err := p.clientDB.NewUpdate().
		Model(&userDAO).
		Where("id = ?", updatedUser.ID).
		OmitZero().
		ExcludeColumn("roles", "id").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return user.User{}, fmt.Errorf("unable to update user: %w", err)
	}

	return userFromDAOToDomain(userDAO), nil
}

// UpdateUserRoles will only update user roles
func (p *pgPersistence) UpdateUserRoles(ctx context.Context, updatedUser user.User) (user.User, error) {
	userDAO := userFromDomainToDAO(updatedUser)

	_, err := p.clientDB.NewUpdate().
		Model(&userDAO).
		Where("id = ?", updatedUser.ID).
		OmitZero().
		Column("roles").
		Returning("*").
		Exec(ctx)
	if err != nil {
		return user.User{}, fmt.Errorf("unable to update user roles: %w", err)
	}

	return userFromDAOToDomain(userDAO), nil
}

func (p *pgPersistence) DeleteUser(ctx context.Context, userIDToDelete user.ID) error {
	_, err := p.clientDB.NewDelete().
		Model((*userDAO)(nil)).
		Where("id = ?", userIDToDelete).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("unable to delete patient id: %d, err: %w", userIDToDelete, err)
	}

	return nil
}
