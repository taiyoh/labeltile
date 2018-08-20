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
