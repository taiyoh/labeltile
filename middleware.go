package labeltile

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/taiyoh/labeltile/app"
)

// UserTokenMiddleware provides framework's whole request filter using UserTokenSerializer
type UserTokenMiddleware struct {
	serializer app.UserTokenSerializer
}

// Execute provides wrapping request by UserTokenMiddleware
func (m *UserTokenMiddleware) Execute(c *gin.Context) {
	claims := m.captureClaims(c.Request.Header.Get("Authorization"))
	if claims != nil {
		if id := claims.FindUserID(); id != "" {
			c.Set("userID", id)
		}
	}
	c.Next()
	if id, ok := c.Get("userID"); ok {
		if claims == nil {
			claims = m.serializer.NewClaims()
		}
		claims.UserID(id.(string))
	}
	if claims != nil {
		c.Writer.Header().Set("Authorization", m.buildNewToken(claims))
	}
}

func (m *UserTokenMiddleware) captureClaims(header string) app.UserTokenClaims {
	auths := strings.SplitN(header, " ", 2)
	if len(auths) != 2 || auths[0] != "Bearer" || auths[1] == "" {
		return nil
	}
	if claims, _ := m.serializer.Deserialize(auths[1]); claims != nil {
		return claims
	}
	return nil
}

func (m *UserTokenMiddleware) buildNewToken(claims app.UserTokenClaims) string {
	token, _ := m.serializer.Serialize(claims)
	authStr := []string{"Bearer"}
	if token != "" {
		authStr = append(authStr, token)
	}
	return strings.Join(authStr, " ")
}

// SetupUserTokenMiddleware provides deserialization token from header and serialization token to header
func SetupUserTokenMiddleware(tokenProvider interface {
	UserTokenSerializer() app.UserTokenSerializer
}) gin.HandlerFunc {
	m := &UserTokenMiddleware{serializer: tokenProvider.UserTokenSerializer()}
	return m.Execute
}
