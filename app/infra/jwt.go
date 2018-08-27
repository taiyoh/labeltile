package infra

import (
	"errors"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
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
func (s *UserTokenSerializer) NewClaims() *UserTokenClaims {
	now := time.Now()
	return &UserTokenClaims{
		claims: map[string]interface{}{
			"iss": "labeltile",
			"iat": strconv.FormatInt(now.Unix(), 10),
			"exp": strconv.FormatInt(now.Add(time.Hour*time.Duration(s.expireHour)).Unix(), 10),
		},
	}
}

// RestoreClaims returns UserTokenClaims object from Claims object
func (s *UserTokenSerializer) RestoreClaims(cl interface{}) *UserTokenClaims {
	claims := cl.(jwt.MapClaims)
	c := map[string]interface{}{}
	for k, v := range claims {
		c[k] = v
	}
	return &UserTokenClaims{claims: c}
}

// Claims returns plain map object
func (c *UserTokenClaims) Claims() map[string]interface{} {
	return c.claims
}

// UserID set user id to claims
func (c *UserTokenClaims) UserID(id string) {
	c.claims["sub"] = id
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
func (s *UserTokenSerializer) Deserialize(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.SecretKey(), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("broken user token")
	}

	claims := s.RestoreClaims(token.Claims)
	if claims.Expired() {
		return nil, errors.New("user token is expired")
	}

	return claims.Claims(), nil
}

// Serialize returns token from plain claims
func (s *UserTokenSerializer) Serialize(claims app.UserTokenClaims) (string, error) {
	token := jwt.New(jwt.GetSigningMethod(s.signingMethod))

	token.Claims = jwt.MapClaims(claims.Claims())

	tokenString, err := token.SignedString(s.SecretKey())
	return tokenString, err
}
