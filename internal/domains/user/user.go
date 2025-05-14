package user

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	RoleAdmin            Role = "admin"
	RoleDoctor           Role = "doctor"
	MinimalMandatoryRole      = RoleDoctor
	// isAdminStandaloneRole at will always validate required roles if admin role is detected
	isAdminStandaloneRole = true
)

type User struct {
	ID       ID       `json:"id"`
	Email    Email    `json:"email"`
	Password Password `json:"-"`
	Roles    Roles    `json:"roles"`
}

type ID int64

type Email string

func (e Email) IsValid() bool {
	return regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(string(e))
}

type Password string

func (p Password) IsValid() bool {
	return len(p) >= 8
}

type Roles []Role

func (r Roles) AreValid() bool {
	foundMandatoryRole := false
	for _, role := range r {
		if !role.IsValid() {
			return false
		}

		if role == MinimalMandatoryRole {
			foundMandatoryRole = true
		}
	}

	return foundMandatoryRole
}

// parse roles from a string
func NewRolesFromRolesString(rolesString string) (r Roles, err error) {
	roles := regexp.MustCompile(",").Split(rolesString, -1)
	for _, role := range roles {
		if availableRoles[Role(role)] {
			r = append(r, Role(role))
		} else {
			err = errors.Join(err, errors.New("invalid role: "+role))
		}
	}

	return r, err
}

func ValidateRequiredRoles(requiredRoles Roles, currentRoles Roles) error {
	var missing Roles

	// isAdminStandaloneRole true will force validation if admin role is present
	// no matter if other roles are missing
	if isAdminStandaloneRole && currentRoles.Has(RoleAdmin) {
		return nil
	}

	for _, reqRole := range requiredRoles {
		if !currentRoles.Has(reqRole) {
			missing = append(missing, reqRole)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required roles: %s", missing.String())
	}

	return nil
}

func (r Roles) Has(role Role) bool {
	for _, existing := range r {
		if existing == role {
			return true
		}
	}
	return false
}

// roles to a single comma sperated string
func (r Roles) String() string {
	parts := make([]string, len(r))
	for i, role := range r {
		parts[i] = role.String()
	}

	return strings.Join(parts, ",")
}

type Role string

func (r Role) String() string {
	return string(r)
}

var availableRoles = map[Role]bool{
	RoleAdmin:  true,
	RoleDoctor: true,
}

func (r Role) IsValid() bool {
	_, exists := availableRoles[r]
	return exists
}
