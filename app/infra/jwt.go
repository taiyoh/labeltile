package infra

import (
	"errors"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// UserTokenSerializer provides JWT based token serialization and deserialization
type UserTokenSerializer struct {
	signingMethod string
	expireHour    uint32
	secretKey     string
}

func isExpired(claims jwt.MapClaims) bool {
	d, ok := claims["expireDate"]
	if !ok {
		return false
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

	var claims jwt.MapClaims
	var ok bool
	claims, ok = token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("broken user token")
	}
	if isExpired(claims) {
		return nil, errors.New("user token is expired")
	}

	c := map[string]interface{}{}
	for k, v := range claims {
		c[k] = v
	}

	return c, nil
}

// Serialize returns token from plain claims
func (s *UserTokenSerializer) Serialize(claims map[string]interface{}) (string, error) {
	token := jwt.New(jwt.GetSigningMethod(s.signingMethod))

	c := jwt.MapClaims{}
	for k, v := range claims {
		c[k] = v
	}
	if _, exists := claims["expireDate"]; !exists {
		t := time.Now().Add(time.Hour * time.Duration(s.expireHour))
		c["expireDate"] = strconv.FormatInt(t.Unix(), 10)
	}

	token.Claims = c

	tokenString, err := token.SignedString(s.SecretKey())
	return tokenString, err
}
