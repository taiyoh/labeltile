package app

import "github.com/taiyoh/labeltile/app/domain"

// Container is interface for infra implementation
type Container interface {
	UserTokenSerializer() UserTokenSerializer
	OAuth2Google() OAuth2Google
	UsepRepository() domain.UserRepository
	RoleRepository() *domain.RoleRepository
}

// UserTokenSerializer is interface for user token serialization
type UserTokenSerializer interface {
	SecretKey() []byte
	Serialize(claims map[string]interface{}) (string, error)
	Deserialize(tokenString string) (map[string]interface{}, error)
}

// OAuth2Google is interface for oauth2 library using google provider
type OAuth2Google interface {
	AuthCodeURL(state string) string
}

// CtxKey is access key for context.Context
type CtxKey string

var (
	// RequestUserCtxKey is access key for requestUser in context
	RequestUserCtxKey = CtxKey("requestUser")
	// ContainerCtxKey is access key for container in context
	ContainerCtxKey = CtxKey("container")
)

// UserDTO is data transfer object for user domain
type UserDTO struct {
	ID    string   `json:"id"`
	Mail  string   `json:"mail"`
	Roles []string `json:"roles"`
}
