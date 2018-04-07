package auth

import (
	"fmt"

	"github.com/taiyoh/labeltile/auth/domain"
)

// MockRegistry is infra aggregation for test
type MockRegistry struct {
	Registry
	userRepo           *MockUserRepository
	userPermissionRepo *MockUserPermissionRepository
}

// NewMockRegistry returns MockRegistry
func NewMockRegistry() Registry {
	return &MockRegistry{
		userRepo: &MockUserRepository{
			userID:   1,
			userPool: []*domain.User{},
		},
		userPermissionRepo: &MockUserPermissionRepository{
			permissionID:     1,
			FoundPermissions: []*domain.UserPermission{},
		},
	}
}

// UserRepository returns object which is implemented with domain's UserRepository interface
func (r *MockRegistry) UserRepository() domain.UserRepository {
	return r.userRepo
}

// UserPermissionRepository returns object which is implemented with domain's UserPermissionRepository interface
func (r *MockRegistry) UserPermissionRepository() domain.UserPermissionRepository {
	return r.userPermissionRepo
}

// MockUserRepository is UserRepository implementation for test
type MockUserRepository struct {
	domain.UserRepository
	userPool []*domain.User
	userID   uint
}

// MockUserPermissionRepository is UserPermissionRepository implementation for test
type MockUserPermissionRepository struct {
	domain.UserPermissionRepository
	FoundPermissions []*domain.UserPermission
	permissionID     uint
}

// DispenseID returns UserID value object
func (r *MockUserRepository) DispenseID() domain.UserID {
	id := r.userID
	r.userID++
	return domain.UserID(fmt.Sprintf("user:%d", id))
}

// Find returns User domain object
func (r *MockUserRepository) Find(id domain.UserID) *domain.User {
	for _, u := range r.userPool {
		if u.ID == id {
			return u
		}
	}
	return nil
}

// Save is emulration for perpetuation to database
func (r *MockUserRepository) Save(u *domain.User) {
	r.userPool = append(r.userPool, u)
}

// DispenseID returns UserPermissionID value object
func (r *MockUserPermissionRepository) DispenseID() domain.UserPermissionID {
	id := r.permissionID
	r.permissionID++
	return domain.UserPermissionID(fmt.Sprintf("user_permission:%d", id))
}

// FindAllByRoles returns list of UserPermission domain objects
func (r *MockUserPermissionRepository) FindAllByRoles(roles []domain.UserRoleID) []*domain.UserPermission {
	return r.FoundPermissions
}
