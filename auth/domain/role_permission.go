package domain

// Role is relation model for user and permission
type Role struct {
	ID          RoleID
	Name        string
	Permissions []*Permission
}

// Permission is model for permission of user's action
type Permission struct {
	ID   PermissionID
	Name string
}

const (
	// RoleViewer is role for view only user. unable to edit anything
	RoleViewer = RoleID(iota)
	// RoleEditor is role for edit data.
	RoleEditor = RoleID(iota)
)

const (
	// PermissionView is permission for view labels.
	PermissionView = PermissionID(iota)
	// PermissionEdit is permission for edit labels.
	PermissionEdit = PermissionID(iota)
)

var (
	roles = map[RoleID]*Role{}
)

func init() {
	permView := &Permission{
		ID:   PermissionView,
		Name: "view labels",
	}
	permEdit := &Permission{
		ID:   PermissionEdit,
		Name: "edit labels",
	}

	roles[RoleViewer] = &Role{
		ID:   RoleViewer,
		Name: "viewer",
		Permissions: []*Permission{
			permView,
		},
	}
	roles[RoleEditor] = &Role{
		ID:   RoleEditor,
		Name: "editor",
		Permissions: []*Permission{
			permView,
			permEdit,
		},
	}
}

// RoleRepository provides interface for fetching Role data
type RoleRepository struct{}

// Find returns Role which is given by id
func (r *RoleRepository) Find(id RoleID) *Role {
	role, _ := roles[id]
	return role
}
