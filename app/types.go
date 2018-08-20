package app

type Container interface {
	UserTokenSerializer() UserTokenSerializer
	OAuth2Google() OAuth2Google
}

type UserTokenSerializer interface {
	SecretKey() []byte
	Serialize(claims map[string]interface{}) (string, error)
	Deserialize(tokenString string) (map[string]interface{}, error)
}

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
