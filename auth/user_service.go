package auth

import "github.com/taiyoh/labeltile/auth/domain"

// UserService provides application service about user in auth context
type UserService struct {
	registry Registry
}

// NewUserService returns UserService
func NewUserService(reg Registry) *UserService {
	return &UserService{registry: reg}
}

// Register provides user registration application service
func (s *UserService) Register(mail, role string) {
	repo := s.registry.UserRepository()
	repo.Save(domain.NewUser(repo.DispenseID(), mail, domain.UserRoleID(role)))
}

// AddRole provides attaching role to user
func (s *UserService) AddRole(id, role string) {
	repo := s.registry.UserRepository()
	u := repo.Find(domain.UserID(id))
	u.AddRole(domain.UserRoleID(role))
	repo.Save(u)
}

// DeleteRole provides detaching role from user
func (s *UserService) DeleteRole(id, role string) {
	repo := s.registry.UserRepository()
	u := repo.Find(domain.UserID(id))
	u.DeleteRole(domain.UserRoleID(role))
	repo.Save(u)
}
