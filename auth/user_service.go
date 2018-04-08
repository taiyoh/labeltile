package auth

import (
	"errors"

	"github.com/taiyoh/labeltile/auth/domain"
)

// UserService provides application service about user in auth context
type UserService struct {
	registry Registry
}

// NewUserService returns UserService
func NewUserService(reg Registry) *UserService {
	return &UserService{registry: reg}
}

// Register provides user registration application service
func (s *UserService) Register(mail, role string) error {
	spec := domain.NewUserSpecification(s.registry.UserRepository(), s.registry.UserPermissionRepository())
	if err := spec.IsSpecifiedToRegisterUser(mail, role); err != nil {
		return err
	}
	repo := s.registry.UserRepository()
	repo.Save(domain.NewUser(repo.DispenseID(), mail, domain.UserRoleID(role)))

	return nil
}

// AddRole provides attaching role to user
func (s *UserService) AddRole(id, role string) error {
	spec := domain.NewUserSpecification(s.registry.UserRepository(), s.registry.UserPermissionRepository())
	if err := spec.IsSpecifiedToEditRole(role); err != nil {
		return err
	}
	repo := s.registry.UserRepository()
	u := repo.Find(domain.UserID(id))
	if u != nil {
		return errors.New("user not found")
	}
	u.AddRole(domain.UserRoleID(role))
	repo.Save(u)

	return nil
}

// DeleteRole provides detaching role from user
func (s *UserService) DeleteRole(id, role string) error {
	spec := domain.NewUserSpecification(s.registry.UserRepository(), s.registry.UserPermissionRepository())
	if err := spec.IsSpecifiedToEditRole(role); err != nil {
		return err
	}
	repo := s.registry.UserRepository()
	u := repo.Find(domain.UserID(id))
	if u != nil {
		return errors.New("user not found")
	}
	u.DeleteRole(domain.UserRoleID(role))
	repo.Save(u)

	return nil
}
