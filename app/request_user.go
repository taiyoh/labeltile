package app

import "time"

// RequestUser is interface for JWT mapped data
type RequestUser interface {
	ID() string
	ExpireDate() time.Time
}

// RequestUserImpl provides RequestUser implementation
type RequestUserImpl struct {
	RequestUser
	id         string
	expireDate int64
}

// NewRequestUser returns RequestUserImpl object
func NewRequestUser(id string, expDate int64) *RequestUserImpl {
	return &RequestUserImpl{
		id:         id,
		expireDate: expDate,
	}
}

// ID returns user ID
func (u *RequestUserImpl) ID() string {
	return u.id
}

// ExpireDate returns JWT expire limitation
func (u *RequestUserImpl) ExpireDate() time.Time {
	return time.Unix(u.expireDate, 0)
}
