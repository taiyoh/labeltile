package mock

import "github.com/taiyoh/labeltile/auth/domain"

type UserRepositoryImpl struct {
	domain.UserRepository
	DispenseIDFunc func() domain.UserID
	Users          map[domain.UserID]*domain.User
}

func (r *UserRepositoryImpl) DispenseID() domain.UserID {
	return r.DispenseIDFunc()
}

func (r *UserRepositoryImpl) Find(id domain.UserID) *domain.User {
	u, _ := r.Users[id]
	return u
}

func (r *UserRepositoryImpl) Save(u *domain.User) {
	r.Users[u.ID] = u
}
