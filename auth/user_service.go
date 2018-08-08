package auth

import (
	"errors"

	"github.com/taiyoh/labeltile/auth/domain"
)

// UserService provides application service about user in auth context
type UserService struct {
	uRepo domain.UserRepository
}

// NewUserService returns UserService
func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{uRepo: repo}
}

// Register provides user registration application service
func (s *UserService) Register(mail string) error {
	spec := domain.NewUserSpecification(s.uRepo)
	if err := spec.SpecifyUserRegistration(mail); err != nil {
		return err
	}
	factory := domain.NewUserFactory(s.uRepo)
	s.uRepo.Save(factory.Build(domain.UserMail(mail)))

	return nil
}

// AddRole provides attaching role to user
func (s *UserService) AddRole(opid, tgtid, role string) error {
	spec := domain.NewUserSpecification(s.registry.UserPermissionRepository())
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
