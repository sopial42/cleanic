package persistence

import (
	user "github.com/sopial42/cleanic/internal/domains/user"
	"github.com/uptrace/bun"
)

type userDAO struct {
	bun.BaseModel `bun:"table:users"`

	ID       int64    `bun:"id,pk,autoincrement"`
	Email    string   `bun:"email,notnull,unique"`
	Password string   `bun:"password,notnull"`
	Roles    []string `bun:"roles,type:jsonb,notnull"`
}

func userFromDomainToDAO(user user.User) userDAO {
	userDAO := userDAO{
		ID:       int64(user.ID),
		Email:    string(user.Email),
		Password: string(user.Password),
	}

	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = string(role)
	}

	userDAO.Roles = roles
	return userDAO
}

func userFromDAOToDomain(userDAO userDAO) user.User {
	domainUser := user.User{
		ID:       user.ID(userDAO.ID),
		Email:    user.Email(userDAO.Email),
		Password: user.Password(userDAO.Password),
	}

	// Convert roles from string slice to domain roles
	domainRoles := make([]user.Role, len(userDAO.Roles))
	for i, role := range userDAO.Roles {
		domainRoles[i] = user.Role(role)
	}

	domainUser.Roles = domainRoles
	return domainUser
}

func userFromDAOsToDomains(userDAOs []userDAO) []user.User {
	users := make([]user.User, len(userDAOs))
	for i, userDAO := range userDAOs {
		users[i] = userFromDAOToDomain(userDAO)
	}

	return users
}
