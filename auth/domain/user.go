package domain

// User is model for accsessing account
type User struct {
	ID    UserID
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

type UserRepository interface {
	Find(id UserID) *User
	Save(u User)
}

type UserPermissionRepository interface {
	FindAllByRoles(roles []UserRoleID) []*UserPermission
}

// NewUser returns initialized user object
func NewUser(id UserID, role UserRoleID) *User {
	return &User{
		ID:    id,
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
