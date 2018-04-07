package domain

type User struct {
	ID    UserID
	Roles []UserRoleID
}

type UserRole struct {
	ID          UserRoleID
	Name        string
	Permissions []UserPermissionID
}

type UserPermission struct {
	ID   UserPermissionID
	Name string
}

func NewUser(id UserID, role UserRoleID) *User {
	return &User{
		ID:    id,
		Roles: []UserRoleID{role},
	}
}

func (u *User) AddRole(r UserRoleID) {
	u.Roles = append(u.Roles, r)
}

func (u *User) DeleteRole(r UserRoleID) {
	rolelist := []UserRoleID{}
	for _, ur := range u.Roles {
		if ur != r {
			rolelist = append(rolelist, ur)
		}
	}
	u.Roles = rolelist
}
