package domain

import (
	"errors"
	"strconv"
)

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
	// RoleManageUser is role for managing users.
	RoleManageUser = RoleID(iota)
)

const (
	// PermissionView is permission for view labels.
	PermissionView = PermissionID(iota)
	// PermissionEdit is permission for edit labels.
	PermissionEdit = PermissionID(iota)
	// PermissionManageUser is permission for managing users.
	PermissionManageUser = PermissionID(iota)
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

	permMngUser := &Permission{
		ID:   PermissionManageUser,
		Name: "manage user",
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
	roles[RoleManageUser] = &Role{
		ID:   RoleManageUser,
		Name: "manager",
		Permissions: []*Permission{
			permMngUser,
		},
	}
}

// RoleRepository provides interface for fetching Role data
type RoleRepository struct{}

// FindMultiByPermission returns Role slices which have given permission
func (r *RoleRepository) FindMultiByPermission(id PermissionID) []*Role {
	roleList := []*Role{}
	for _, role := range roles {
		for _, p := range role.Permissions {
			if p.ID == id {
				roleList = append(roleList, role)
			}
		}
	}
	return roleList
}

// RoleSpecification provides validation methods for role modification
type RoleSpecification struct {
	rRepo *RoleRepository
}

// NewRoleSpecification returns RoleSpecification struct
func NewRoleSpecification(r *RoleRepository) *RoleSpecification {
	return &RoleSpecification{rRepo: r}
}

// SpecifyEditRole returns whether operator is editable or not
func (s *RoleSpecification) SpecifyEditRole(op *User, roleIDList []string) error {
	canOperates := map[RoleID]struct{}{}
	for _, role := range s.rRepo.FindMultiByPermission(PermissionManageUser) {
		canOperates[role.ID] = struct{}{}
	}
	var canOperate bool
	for _, r := range op.Roles {
		if _, ok := canOperates[r]; ok {
			canOperate = true
			break
		}
	}
	if !canOperate {
		return errors.New("not permitted")
	}

	for _, rid := range roleIDList {
		if irid, err := strconv.Atoi(rid); err == nil {
			if _, ok := roles[RoleID(irid)]; !ok {
				return errors.New("role not found")
			}
		} else {
			return errors.New("role not found")
		}
	}

	return nil
}
