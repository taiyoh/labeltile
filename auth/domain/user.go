package domain

type userRoles []UserRoleID

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
}

// UserFactory is builder for User
type UserFactory struct {
	uRepo UserRepository
}

// NewUserFactory returns UserFactory struct
func NewUserFactory(r UserRepository) *UserFactory {
	return &UserFactory{
		uRepo: r,
	}
}

// Build returns User struct
func (f *UserFactory) Build(m UserMail) *User {
	id := f.uRepo.DispenseID()
	return &User{
		ID:    id,
		Mail:  m,
		Roles: userRoles{UserRoleViewer},
	}
}

func (r userRoles) Add(id UserRoleID) userRoles {
	nr := r[:]
	return append(nr, id)
}

func (r userRoles) Delete(id UserRoleID) userRoles {
	nr := userRoles{}
	for _, ro := range r {
		if ro != id {
			nr = append(nr, ro)
		}
	}
	return nr
}

// AddRole set role to user
func (u *User) AddRole(r UserRoleID) *User {
	return &User{
		ID:    u.ID,
		Mail:  u.Mail,
		Roles: u.Roles.Add(r),
	}
}

// DeleteRole unset role from user
func (u *User) DeleteRole(r UserRoleID) *User {
	return &User{
		ID:    u.ID,
		Mail:  u.Mail,
		Roles: u.Roles.Delete(r),
	}
}
