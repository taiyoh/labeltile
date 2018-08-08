package domain

import (
	"errors"
	"net/mail"
)

type userRoles []RoleID

// User is model for accsessing account
type User struct {
	ID    UserID
	Mail  UserMail
	Roles userRoles
}

// UserRepository is interface for fetching User aggregation from perpetuation layer
type UserRepository interface {
	DispenseID() UserID
	Find(id UserID) *User
	Save(u *User)
	FindByMail(m string) *User
}

// UserFactory is builder for User
type UserFactory struct {
	uRepo UserRepository
}

// UserSpecification provides validation methods for user registration
type UserSpecification struct {
	uRepo UserRepository
}

// NewUserFactory returns UserFactory struct
func NewUserFactory(r UserRepository) *UserFactory {
	return &UserFactory{
		uRepo: r,
	}
}

// NewUserSpecification returns UserSpecification struct
func NewUserSpecification(r UserRepository) *UserSpecification {
	return &UserSpecification{
		uRepo: r,
	}
}

// Build returns User struct
func (f *UserFactory) Build(m UserMail) *User {
	id := f.uRepo.DispenseID()
	return &User{
		ID:    id,
		Mail:  m,
		Roles: userRoles{RoleViewer},
	}
}

// SpecifyUserRegistration returns whether enable to register user or not
func (s *UserSpecification) SpecifyUserRegistration(addr string) error {
	e, err := mail.ParseAddress(addr)
	if err != nil {
		return err
	}

	if e.Address != addr {
		return errors.New("conatins something other than E-mail address")
	}

	u := s.uRepo.FindByMail(e.Address)
	if u != nil {
		return errors.New("already registered")
	}
	return nil
}

func (r userRoles) Add(id RoleID) userRoles {
	nr := userRoles{}
	var roleExists bool
	for _, ro := range r {
		if ro == id {
			roleExists = true
		}
		nr = append(nr, ro)
	}
	if roleExists {
		return nr
	}
	return append(nr, id)
}

func (r userRoles) Delete(id RoleID) userRoles {
	nr := userRoles{}
	for _, ro := range r {
		if ro != id {
			nr = append(nr, ro)
		}
	}
	return nr
}

// AddRole set role to user
func (u *User) AddRole(r RoleID) *User {
	return &User{
		ID:    u.ID,
		Mail:  u.Mail,
		Roles: u.Roles.Add(r),
	}
}

// DeleteRole unset role from user
func (u *User) DeleteRole(r RoleID) *User {
	return &User{
		ID:    u.ID,
		Mail:  u.Mail,
		Roles: u.Roles.Delete(r),
	}
}
