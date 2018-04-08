package domain

import "errors"

// User is model for accsessing account
type User struct {
	ID    UserID
	Mail  string
	Roles []UserRoleID
}

// UserRole is relation model for user and permission
type UserRole struct {
	ID          UserRoleID
	Name        string
	Permissions []UserPermissionID
}

// UserPermission is model for permission of user's action
type UserPermission struct {
	ID   UserPermissionID
	Name string
}

// UserRepository is interface for fetching User aggregation from perpetuation layer
type UserRepository interface {
	DispenseID() UserID
	Find(id UserID) *User
	Save(u *User)
}

// UserPermissionRepository is interface for fetching UserPermission aggregation from perpetuation layer
type UserPermissionRepository interface {
	DispenseID() UserPermissionID
	FindAllByRoles(roles []UserRoleID) []*UserPermission
}

// UserSpecification provides specification and validation for user
type UserSpecification struct {
	userRepo       UserRepository
	permissionRepo UserPermissionRepository
}

// NewUserSpecification returns UserSpecification object
func NewUserSpecification(uRepo UserRepository, pRepo UserPermissionRepository) *UserSpecification {
	return &UserSpecification{
		userRepo:       uRepo,
		permissionRepo: pRepo,
	}
}

// IsSpecifiedToRegisterUser provides validation for registering user
func (s *UserSpecification) IsSpecifiedToRegisterUser(mail, role string) error {
	return s.isValidRole(role)
}

// IsSpecifiedToEditRole provides vaildation for adding or deleting user's role
func (s *UserSpecification) IsSpecifiedToEditRole(role string) error {
	return s.isValidRole(role)
}

func (s *UserSpecification) isValidRole(role string) error {
	perms := s.permissionRepo.FindAllByRoles([]UserRoleID{UserRoleID(role)})
	if len(perms) == 0 {
		return errors.New("role not registered")
	}
	return nil
}

// NewUser returns initialized user object
func NewUser(id UserID, mail string, role UserRoleID) *User {
	return &User{
		ID:    id,
		Mail:  mail,
		Roles: []UserRoleID{role},
	}
}

// AddRole set role to user
func (u *User) AddRole(r UserRoleID) {
	u.Roles = append(u.Roles, r)
}

// DeleteRole unset role from user
func (u *User) DeleteRole(r UserRoleID) {
	rolelist := []UserRoleID{}
	for _, ur := range u.Roles {
		if ur != r {
			rolelist = append(rolelist, ur)
		}
	}
	u.Roles = rolelist
}
