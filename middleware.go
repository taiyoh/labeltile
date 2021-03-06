package labeltile

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/taiyoh/labeltile/app"
)

const (
	// CtxSessionIDKey is for handling session id in Context
	CtxSessionIDKey = "sessionID"
)

// UserTokenMiddleware provides framework's whole request filter using UserTokenSerializer
type UserTokenMiddleware struct {
	serializer app.UserTokenSerializer
}

// Execute provides wrapping request by UserTokenMiddleware
func (m *UserTokenMiddleware) Execute(c *gin.Context) {
	claims := m.captureClaims(c.Request.Header.Get("Authorization"))
	if id := claims.FindUserID(); id != "" {
		c.Set("userID", id)
	}
	c.Set(CtxSessionIDKey, claims.FindSessionID())
	c.Next()
	if id, ok := c.Get("userID"); ok {
		claims.UserID(id.(string))
	}
	c.Writer.Header().Set("Authorization", m.buildNewToken(claims))
}

func (m *UserTokenMiddleware) captureClaims(header string) app.UserTokenClaims {
	auths := strings.SplitN(header, " ", 2)
	if len(auths) != 2 || auths[0] != "Bearer" || auths[1] == "" {
		return m.serializer.NewClaims()
	}
	if claims, _ := m.serializer.Deserialize(auths[1]); claims != nil {
		return claims
	}
	return m.serializer.NewClaims()
}

func (m *UserTokenMiddleware) buildNewToken(claims app.UserTokenClaims) string {
	token, _ := m.serializer.Serialize(claims)
	return strings.Join([]string{"Bearer", token}, " ")
}

// SessionMiddleware provides framework's whole request filter using SessionStorage
type SessionMiddleware struct {
	session app.SessionStorage
}

// Execute provides wrapping request by SessionMiddleware
func (m *SessionMiddleware) Execute(c *gin.Context) {
	sessionID, exists := c.Get(CtxSessionIDKey)
	if !exists {
		c.Next()
		return
	}
	sid, _ := sessionID.(string)
	session := m.session.Find(sid)
	if session == nil {
		session = m.session.New(sid)
	}
	c.Set("session", session)
	c.Next()
	m.session.Save(sid, session)
}

// SetupUserTokenMiddleware provides deserialization token from header and serialization token to header
func SetupUserTokenMiddleware(tokenProvider interface {
	UserTokenSerializer() app.UserTokenSerializer
}) gin.HandlerFunc {
	m := &UserTokenMiddleware{serializer: tokenProvider.UserTokenSerializer()}
	return m.Execute
}

// SetupSessionMiddleware provides session storage operation request by request
func SetupSessionMiddleware(sessionProvider interface {
	SessionStorage() app.SessionStorage
}) gin.HandlerFunc {
	m := &SessionMiddleware{session: sessionProvider.SessionStorage()}
	return m.Execute
}
