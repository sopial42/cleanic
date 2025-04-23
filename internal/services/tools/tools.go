package tools

import (
	"fmt"

	"github.com/sopial42/cleanic/internal/domains/user"
)

func EnsureUserOwnershipOrRoleAdmin(reqUser user.User, userIDToUpdate user.ID) error {
	if reqUser.ID != userIDToUpdate && !reqUser.Roles.Has(user.RoleAdmin) {
		return fmt.Errorf("unauthorized to update another user")
	}

	return nil
}
