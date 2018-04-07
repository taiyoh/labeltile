package auth

import "github.com/taiyoh/labeltile/auth/domain"

// MockRegistry is infra aggregation for test
type MockRegistry struct {
	Registry
	userRepo           *MockUserRepository
	userPermissionRepo *MockUserPermissionRepository
}

// NewMockRegistry returns MockRegistry
func NewMockRegistry() Registry {
	return &MockRegistry{
		userRepo: &MockUserRepository{},
		userPermissionRepo: &MockUserPermissionRepository{
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
	FoundUser *domain.User
	SavedUser *domain.User
}

// MockUserPermissionRepository is UserPermissionRepository implementation for test
type MockUserPermissionRepository struct {
	domain.UserPermissionRepository
	FoundPermissions []*domain.UserPermission
}

// Find returns User domain object
func (r *MockUserRepository) Find(id domain.UserID) *domain.User {
	return r.FoundUser
}

// Save is emulration for perpetuation to database
func (r *MockUserRepository) Save(u *domain.User) {
	r.SavedUser = u
}

// FindAllByRoles returns list of UserPermission domain objects
func (r *MockUserPermissionRepository) FindAllByRoles(roles []domain.UserRoleID) []*domain.UserPermission {
	return r.FoundPermissions
}
