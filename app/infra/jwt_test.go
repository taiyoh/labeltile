package infra_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/taiyoh/labeltile/app/infra"
)

func TestSerialize(t *testing.T) {
	s := infra.NewUserTokenSerializer("HS512", "foobar", 1)

	now := time.Now()
	claims := map[string]interface{}{
		"foo": "bar",
	}
	var tokenString string
	var err error
	tokenString, err = s.Serialize(claims)
	if tokenString == "" {
		t.Error("failed to serialize token")
	}
	if err != nil {
		t.Error("error found:", err)
	}

	var c map[string]interface{}
	c, err = s.Deserialize(tokenString)
	if err != nil {
		t.Error("error found: " + err.Error())
	} else if v, exists := c["foo"]; !exists {
		t.Error("not found 'foo' in deserialized claims")
	} else if vs, ok := v.(string); !ok {
		t.Error("'foo' is not string")
	} else if vs != "bar" {
		t.Error("contained value is wrong")
	}

	if c["expireDate"].(string) != strconv.FormatInt(now.Add(time.Hour).Unix(), 10) {
		t.Error("expireDate is wrong")
	}

	if _, err := s.Deserialize(tokenString + "foobar"); err == nil {
		t.Error("user token has broken")
	}
}

func TestExpiredDateToken(t *testing.T) {
	s := infra.NewUserTokenSerializer("HS512", "foobar", 1)

	now := time.Now()
	claims := map[string]interface{}{
		"foo": "bar",
	}

	claims["expireDate"] = strconv.FormatInt(now.Unix(), 10)
	tokenString, _ := s.Serialize(claims)

	if _, err := s.Deserialize(tokenString); err == nil {
		t.Error("user token has expired")
	}

	claims["expireDate"] = strconv.FormatInt(now.Add(-time.Hour).Unix(), 10)
	tokenString, _ = s.Serialize(claims)

	if _, err := s.Deserialize(tokenString); err == nil {
		t.Error("user token has expired")
	}

}
