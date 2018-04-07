package auth

import "github.com/taiyoh/labeltile/auth/domain"

// Registry is infra aggregation interface for auth context
type Registry interface {
	UserRepository() domain.UserRepository
	UserPermissionRepository() domain.UserPermissionRepository
}
