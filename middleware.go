package labeltile

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/taiyoh/labeltile/app"
)

func captureClaims(header string, s app.UserTokenSerializer) map[string]interface{} {
	auths := strings.SplitN(header, " ", 2)
	if len(auths) != 2 || auths[0] != "Bearer" || auths[1] == "" {
		return map[string]interface{}{}
	}
	if claims, _ := s.Deserialize(auths[1]); claims != nil {
		return claims
	}
	return map[string]interface{}{}
}

// UserTokenMiddleware provides deserialization token from header and serialization token to header
func UserTokenMiddleware(tokenProvider interface {
	UserTokenSerializer() app.UserTokenSerializer
}) gin.HandlerFunc {
	s := tokenProvider.UserTokenSerializer()
	return func(c *gin.Context) {
		claims := captureClaims(c.Request.Header.Get("Authorization"), s)
		if id, ok := claims["userID"]; ok {
			c.Set("userID", id.(string))
		}
		c.Next()
		if id, ok := c.Get("userID"); ok {
			claims["userID"] = id.(string)
		}
		token, _ := s.Serialize(claims)
		joinStr := []string{"Bearer"}
		if token != "" {
			joinStr = append(joinStr, token)
		}
		c.Writer.Header().Set("Authorization", strings.Join(joinStr, " "))
	}
}
