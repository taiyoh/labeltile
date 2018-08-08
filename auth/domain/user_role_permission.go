package domain

// UserRole is relation model for user and permission
type UserRole struct {
	ID          UserRoleID
	Name        string
	Permissions []*UserPermission
}

// UserPermission is model for permission of user's action
type UserPermission struct {
	ID   UserPermissionID
	Name string
}

const (
	// UserRoleViewer is role for view only user. unable to edit anything
	UserRoleViewer = UserRoleID(iota)
	// UserRoleEditor is role for edit data.
	UserRoleEditor = UserRoleID(iota)
)

const (
	// UserPermissionView is permission for view labels.
	UserPermissionView = UserPermissionID(iota)
	// UserPermissionEdit is permission for edit labels.
	UserPermissionEdit = UserPermissionID(iota)
)

var (
	roles = map[UserRoleID]*UserRole{}
)

func init() {
	permView := &UserPermission{
		ID:   UserPermissionView,
		Name: "view labels",
	}
	permEdit := &UserPermission{
		ID:   UserPermissionEdit,
		Name: "edit labels",
	}

	roles[UserRoleViewer] = &UserRole{
		ID:   UserRoleViewer,
		Name: "viewer",
		Permissions: []*UserPermission{
			permView,
		},
	}
	roles[UserRoleEditor] = &UserRole{
		ID:   UserRoleEditor,
		Name: "editor",
		Permissions: []*UserPermission{
			permView,
			permEdit,
		},
	}
}

type UserRoleRepository struct{}

func (r *UserRoleRepository) Find(id UserRoleID) *UserRole {
	role, _ := roles[id]
	return role
}
