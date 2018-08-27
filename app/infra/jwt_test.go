package infra_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/taiyoh/labeltile/app/infra"
)

func TestSerialize(t *testing.T) {
	s := infra.NewUserTokenSerializer("HS512", "foobar", 1)

	claims := s.NewClaims()
	claims.UserID("hogefuga")
	var tokenString string
	var err error
	tokenString, err = s.Serialize(claims)
	if tokenString == "" {
		t.Error("failed to serialize token")
	}
	if err != nil {
		t.Error("error found:", err)
	}

	now := time.Now()
	if c, err := s.Deserialize(tokenString); err != nil {
		t.Error("error found: " + err.Error())
	} else {
		if v, exists := c["iat"]; !exists {
			t.Error("not found 'iat' in deserialized claims")
		} else if vs, ok := v.(string); !ok {
			t.Error("'iat' is not string")
		} else if vs != strconv.FormatInt(now.Unix(), 10) {
			t.Error("'iat' value is wrong")
		}
		if v, exists := c["sub"]; !exists {
			t.Error("not found 'sub' in deserialized claims")
		} else if vs, ok := v.(string); !ok {
			t.Error("'sub' is not string")
		} else if vs != "hogefuga" {
			t.Error("'sub' value is wrong")
		}

		if c["exp"].(string) != strconv.FormatInt(now.Add(time.Hour).Unix(), 10) {
			t.Error("exp is wrong")
		}
	}

	if _, err := s.Deserialize(tokenString + "foobar"); err == nil {
		t.Error("user token has broken")
	}
}

func TestExpiredDateToken(t *testing.T) {
	s := infra.NewUserTokenSerializer("HS512", "foobar", 1)

	now := time.Now()
	claims := s.NewClaims()
	innerClaims := claims.Claims()

	innerClaims["exp"] = strconv.FormatInt(now.Unix(), 10)
	claims = s.RestoreClaims(jwt.MapClaims(innerClaims))
	tokenString, _ := s.Serialize(claims)

	if _, err := s.Deserialize(tokenString); err == nil {
		t.Error("user token has expired")
	}

	innerClaims["exp"] = strconv.FormatInt(now.Add(-time.Hour).Unix(), 10)
	claims = s.RestoreClaims(jwt.MapClaims(innerClaims))
	tokenString, _ = s.Serialize(claims)

	if _, err := s.Deserialize(tokenString); err == nil {
		t.Error("user token has expired")
	}

	delete(innerClaims, "exp")
	claims = s.RestoreClaims(jwt.MapClaims(innerClaims))
	if !claims.Expired() {
		t.Error("claims is expired when exp key is not exists")
	}
}
