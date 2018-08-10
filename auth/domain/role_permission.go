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

// ConvertToID returns RoleID when specified id has given.
func (r *RoleRepository) ConvertToID(id string) (*RoleID, error) {
	if irid, err := strconv.Atoi(id); err == nil {
		rid := RoleID(irid)
		if _, ok := roles[rid]; !ok {
			return nil, errors.New("role not found")
		}
		return &rid, nil
	}
	return nil, errors.New("role not found")
}

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

func (s *RoleSpecification) specifyEditRole(op *User, roleIDs []RoleID) error {
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

	for _, rid := range roleIDs {
		if rid == RoleViewer {
			return errors.New("cannot edit Viewer Role")
		}
	}

	return nil
}

// ConvertRoleToID returns RoleID list given from role string list
func (s *RoleSpecification) ConvertRoleToID(roles []string) ([]RoleID, error) {
	roleIDs := []RoleID{}

	if len(roles) == 0 {
		return roleIDs, errors.New("require role list")
	}
	for _, r := range roles {
		if rid, err := s.rRepo.ConvertToID(r); err == nil {
			roleIDs = append(roleIDs, *rid)
		}
	}

	if len(roles) != len(roleIDs) {
		return roleIDs, errors.New("invalid role exists")
	}

	return roleIDs, nil
}

// SpecifyAddRole provides whether operator can add role to target or not
func (s *RoleSpecification) SpecifyAddRole(op, tgt *User, roleIDs []RoleID) error {
	return s.specifyEditRole(op, roleIDs)
}

// SpecifyDeleteRole provides whether operator can delete role from target or not
func (s *RoleSpecification) SpecifyDeleteRole(op, tgt *User, roleIDs []RoleID) error {
	if err := s.specifyEditRole(op, roleIDs); err != nil {
		return err
	}

	for _, rid := range roleIDs {
		if rid == RoleManageUser && op.ID == tgt.ID {
			return errors.New("cannot detach self Manager role")
		}
	}

	return nil
}

// SpecifyRegisterUser returns whether operator has Manager role or not.
func (s *RoleSpecification) SpecifyRegisterUser(op *User) error {
	return s.specifyEditRole(op, []RoleID{RoleManageUser})
}
