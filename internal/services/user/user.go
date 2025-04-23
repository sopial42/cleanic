package user

import (
	"context"
	"fmt"

	user "github.com/sopial42/cleanic/internal/domains/user"
	"github.com/sopial42/cleanic/internal/services/tools"
	"golang.org/x/crypto/bcrypt"
)

type userSVC struct {
	persistence Persistence
}

func NewUserService(persistence Persistence) Service {
	return &userSVC{
		persistence: persistence,
	}
}

func (u *userSVC) Register(ctx context.Context, newUser user.User) (user.User, error) {
	// Assign default role
	newUser.Roles = []user.Role{user.MinimalMandatoryRole}

	// Hash password
	hash, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	newUser.Password = string(hash)

	// Check email is correct
	if !newUser.Email.IsValid() {
		return user.User{}, fmt.Errorf("invalid email input: %v", newUser.Email)
	}

	userCreated, err := u.persistence.Insert(ctx, newUser)
	if err != nil {
		return user.User{}, err
	}

	return userCreated, nil
}

func (u *userSVC) Login(ctx context.Context, loginUser user.User) (user.User, SignedJWT, error) {
	// Check email is correct
	if !loginUser.Email.IsValid() {
		return user.User{}, SignedJWT{}, fmt.Errorf("invalid email input: %v", loginUser.Email)
	}
	// Get user by email
	userFound, err := u.persistence.GetUserByEmail(ctx, loginUser.Email)
	if err != nil {
		return user.User{}, SignedJWT{}, err
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(loginUser.Password)); err != nil {
		return user.User{}, SignedJWT{}, fmt.Errorf("invalid password: %w", err)
	}

	jwt, err := generateJWT(userFound.ID, userFound.Roles)
	if err != nil {
		return user.User{}, SignedJWT{}, fmt.Errorf("unable to generate token: %w", err)
	}

	return userFound, jwt, nil
}

func (u *userSVC) GetUsers(ctx context.Context) ([]user.User, error) {
	users, err := u.persistence.ListUsers(ctx)
	if err != nil {
		return []user.User{}, err
	}

	return users, nil
}

func (u *userSVC) GetUserByID(ctx context.Context, id user.ID) (user.User, error) {
	userFound, err := u.persistence.GetUserByID(ctx, id)
	if err != nil {
		return user.User{}, err
	}

	return userFound, nil
}

func (u *userSVC) UpdateUser(ctx context.Context, reqUserID user.ID, newUser user.User) (user.User, error) {
	reqUser, err := u.persistence.GetUserByID(ctx, reqUserID)
	if err != nil {
		return user.User{}, fmt.Errorf("unable to get request's user details: %w", err)
	}

	// Ensure no roles will be updated
	if len(newUser.Roles) > 0 {
		return user.User{}, fmt.Errorf("unauthorize to update user roles")
	}

	// Ensure one field or more will be updated
	if newUser.Email == "" && newUser.Password == "" {
		return user.User{}, fmt.Errorf("unable to update nothing on the user")
	}

	// Ensure user try to update its own user OR is admin
	if err = tools.EnsureUserOwnershipOrRoleAdmin(reqUser, newUser.ID); err != nil {
		return user.User{}, err
	}

	if newUser.Email != "" {
		if !newUser.Email.IsValid() {
			return user.User{}, fmt.Errorf("unable to update user, invalid email: %v", newUser.Email)
		}
	}

	if newUser.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		newUser.Password = string(hash)
	}

	userUpdated, err := u.persistence.UpdateUser(ctx, newUser)
	if err != nil {
		return user.User{}, fmt.Errorf("unable to update user: %w", err)
	}

	return userUpdated, nil
}

func (u *userSVC) UpdateUserRoles(ctx context.Context, reqUserID user.ID, updatedUser user.User) (user.User, error) {
	reqUser, err := u.persistence.GetUserByID(ctx, reqUserID)
	if err != nil {
		return user.User{}, fmt.Errorf("unable to get request's user details: %w", err)
	}

	if !reqUser.Roles.Has(user.RoleAdmin) {
		return user.User{}, fmt.Errorf("unauthorized to updated roles")
	}

	if !updatedUser.Roles.AreValid() {
		return user.User{}, fmt.Errorf("unable to update roles, invalid input: %v", updatedUser.Roles)
	}

	userUpdated, err := u.persistence.UpdateUserRoles(ctx, updatedUser)
	if err != nil {
		return user.User{}, fmt.Errorf("unable to update user: %w", err)
	}

	return userUpdated, nil
}

func (u *userSVC) DeleteUser(ctx context.Context, reqUserID user.ID, userIDToDelete user.ID) error {
	reqUser, err := u.persistence.GetUserByID(ctx, reqUserID)
	if err != nil {
		return fmt.Errorf("unable to get request's user details: %w", err)
	}

	if err = tools.EnsureUserOwnershipOrRoleAdmin(reqUser, userIDToDelete); err != nil {
		return err
	}

	err = u.persistence.DeleteUser(ctx, userIDToDelete)
	if err != nil {
		return fmt.Errorf("unable to update user: %w", err)
	}

	return nil
}
