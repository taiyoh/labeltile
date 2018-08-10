package mock

import "github.com/taiyoh/labeltile/auth/domain"

// UserRepositoryImpl provides mock for UserRepository implementation
type UserRepositoryImpl struct {
	domain.UserRepository
	DispenseIDFunc func() domain.UserID
	Users          map[domain.UserID]*domain.User
}

// LoadUserRepoImpl reutrns UserRepositoryImpl struct
func LoadUserRepoImpl(f func() domain.UserID) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		DispenseIDFunc: f,
		Users:          map[domain.UserID]*domain.User{},
	}
}

// DispenseID is implementation for mock
func (r *UserRepositoryImpl) DispenseID() domain.UserID {
	return r.DispenseIDFunc()
}

// Find is implementation for mock
func (r *UserRepositoryImpl) Find(id string) *domain.User {
	u, _ := r.Users[domain.UserID(id)]
	return u
}

// Save is implementation for mock
func (r *UserRepositoryImpl) Save(u *domain.User) {
	r.Users[u.ID] = u
}

// FindByMail is implementation for mock
func (r *UserRepositoryImpl) FindByMail(addr string) *domain.User {
	for _, u := range r.Users {
		if string(u.Mail) == addr {
			return u
		}
	}
	return nil
}

func (r *UserRepositoryImpl) FindMulti(ids []string) []*domain.User {
	users := []*domain.User{}
	for _, u := range r.Users {
		for _, id := range ids {
			if u.ID == domain.UserID(id) {
				users = append(users, u)
			}
		}
	}
	return users
}
