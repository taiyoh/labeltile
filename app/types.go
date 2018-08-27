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

type UserTokenClaims interface {
	Claims() map[string]interface{}
}

// UserTokenSerializer is interface for user token serialization
type UserTokenSerializer interface {
	SecretKey() []byte
	Serialize(UserTokenClaims) (string, error)
	Deserialize(string) (map[string]interface{}, error)
}

// OAuth2Google is interface for oauth2 library using google provider
type OAuth2Google interface {
	AuthCodeURL(state string) string
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
