package app

import "github.com/taiyoh/labeltile/app/domain"

// Container is interface for infra implementation
type Container interface {
	Register(string, interface{})
	UserTokenSerializer() UserTokenSerializer
	OAuth2Google() OAuth2Google
	UserRepository() domain.UserRepository
	RoleRepository() *domain.RoleRepository
}

// UserTokenClaims is interface for user token claims
type UserTokenClaims interface {
	Claims() map[string]interface{}
	FindUserID() string
	UserID(string)
	Expired() bool
	FindSessionID() string
}

// UserTokenSerializer is interface for user token serialization
type UserTokenSerializer interface {
	SecretKey() []byte
	Serialize(UserTokenClaims) (string, error)
	Deserialize(string) (UserTokenClaims, error)
	NewClaims() UserTokenClaims
	RestoreClaims(map[string]interface{}) UserTokenClaims
}

// OAuth2Google is interface for oauth2 library using google provider
type OAuth2Google interface {
	AuthCodeURL(state string) string
}

// SessionData is interface for session data storage interface for request
type SessionData interface {
	Get(string) interface{}
	Set(string, interface{})
	Remove(string) bool
}

// SessionStorage is interface for Session data operation
type SessionStorage interface {
	Find(string) SessionData
	New(string) SessionData
	Save(string, SessionData)
	Remove(string) bool
}

// CtxKey is access key for context.Context
type CtxKey string

var (
	// UserIDCtxKey is access key for requestUser in context
	UserIDCtxKey = CtxKey("userID")
)

// UserDTO is data transfer object for user domain
type UserDTO struct {
	ID    string   `json:"id"`
	Mail  string   `json:"mail"`
	Roles []string `json:"roles"`
}
