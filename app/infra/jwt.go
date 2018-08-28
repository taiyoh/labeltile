package infra

import (
	"errors"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/taiyoh/labeltile/app"
)

// UserTokenSerializer provides JWT based token serialization and deserialization
type UserTokenSerializer struct {
	app.UserTokenSerializer
	signingMethod string
	expireHour    uint32
	secretKey     string
}

// UserTokenClaims is wrapper for JWT Claims
type UserTokenClaims struct {
	app.UserTokenClaims
	claims map[string]interface{}
}

// NewClaims returns UserTokenClaims object
func (s *UserTokenSerializer) NewClaims() app.UserTokenClaims {
	now := time.Now()
	return &UserTokenClaims{
		claims: map[string]interface{}{
			"iss": "labeltile",
			"iat": strconv.FormatInt(now.Unix(), 10),
			"exp": strconv.FormatInt(now.Add(time.Hour*time.Duration(s.expireHour)).Unix(), 10),
			"sid": uuid.New().String(),
		},
	}
}

// RestoreClaims returns UserTokenClaims object from Claims object
func (s *UserTokenSerializer) RestoreClaims(claims map[string]interface{}) app.UserTokenClaims {
	return &UserTokenClaims{claims: claims}
}

// Claims returns plain map object
func (c *UserTokenClaims) Claims() map[string]interface{} {
	return c.claims
}

// UserID provides setting user id to claims
func (c *UserTokenClaims) UserID(id string) {
	c.claims["sub"] = id
}

// FindUserID returns user id from claims
func (c *UserTokenClaims) FindUserID() string {
	if id, ok := c.claims["sub"]; ok {
		return id.(string)
	}
	return ""
}

// Expired returns whether this claims' expirenation is over or not
func (c *UserTokenClaims) Expired() bool {
	d, ok := c.claims["exp"]
	if !ok {
		return true
	}
	now := time.Now()
	ds, _ := d.(string)
	dt, _ := strconv.Atoi(ds)
	t := time.Unix(int64(dt), 0)
	return now.Equal(t) || now.After(t)
}

// FindSessionID returns session id which is already created
func (c *UserTokenClaims) FindSessionID() string {
	sid, _ := c.claims["sid"].(string)
	return sid
}

// NewUserTokenSerializer returns UserTokenSerializer object
func NewUserTokenSerializer(method, skey string, h uint32) *UserTokenSerializer {
	return &UserTokenSerializer{
		signingMethod: method,
		expireHour:    h,
		secretKey:     skey,
	}
}

// SecretKey returns binarized secret key
func (s *UserTokenSerializer) SecretKey() []byte {
	return []byte(s.secretKey)
}

// Deserialize returns plain claims from token
func (s *UserTokenSerializer) Deserialize(tokenString string) (app.UserTokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.SecretKey(), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("broken user token")
	}

	claims := s.RestoreClaims(map[string]interface{}(token.Claims.(jwt.MapClaims)))
	if claims.Expired() {
		return nil, errors.New("user token is expired")
	}

	return claims, nil
}

// Serialize returns token from plain claims
func (s *UserTokenSerializer) Serialize(claims app.UserTokenClaims) (string, error) {
	token := jwt.New(jwt.GetSigningMethod(s.signingMethod))

	token.Claims = jwt.MapClaims(claims.Claims())

	tokenString, err := token.SignedString(s.SecretKey())
	return tokenString, err
}
