package labeltile_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/taiyoh/labeltile"
	"github.com/taiyoh/labeltile/app/infra"
	"github.com/taiyoh/labeltile/app/infra/mock"

	"github.com/gin-gonic/gin"
)

func TestTokenMiddleware(t *testing.T) {
	router := gin.Default()
	c := mock.LoadContainer()
	s := infra.NewUserTokenSerializer("HS512", "foobar", 1)
	c.Register("UserTokenSerializer", s)
	router.Use(labeltile.UserTokenMiddleware(c))
	router.GET("/", func(c *gin.Context) {
		userID, _ := c.Get("userID")
		c.JSON(http.StatusOK, gin.H{"userID": userID})
	})
	router.GET("/login", func(c *gin.Context) {
		c.Set("userID", "foo")
		c.JSON(http.StatusOK, gin.H{"userID": "foo"})
	})
	t.Run("no Authorization header request", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		auths := strings.SplitN(w.Header().Get("Authorization"), " ", 2)
		if len(auths) != 2 {
			t.Error("wrong Authorization header")
		}
		if auths[0] != "Bearer" {
			t.Error("Authorization header prefix is Bearer")
		}
		claims, err := s.Deserialize(auths[1])
		if err != nil {
			t.Error("deserialize failed:", err)
		}
		if _, ok := claims["userID"]; ok {
			t.Error("ghost userID filled")
		}
		if len(claims) != 1 {
			t.Error("expireDate key only")
		}
	})
	t.Run("no Authorization header request to login", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/login", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		auths := strings.SplitN(w.Header().Get("Authorization"), " ", 2)
		claims, err := s.Deserialize(auths[1])
		if err != nil {
			t.Error("deserialize failed:", err)
		}
		id, userIDexists := claims["userID"]
		if !userIDexists {
			t.Error("userID not found")
		}
		if id.(string) != "foo" {
			t.Error("wrong userID")
		}
	})
	t.Run("broken Authorization header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		brokenJwt := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		req.Header.Set("Authorization", "Bearer "+brokenJwt)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		auths := strings.SplitN(w.Header().Get("Authorization"), " ", 2)
		if auths[1] == brokenJwt {
			t.Error("broken JWT through")
		}

		claims, err := s.Deserialize(auths[1])
		if err != nil {
			t.Error("deserialize failed:", err)
		}
		if _, userIDexists := claims["userID"]; userIDexists {
			t.Error("ghost userID found")
		}
	})
	t.Run("valid Authorization header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		validJwt, _ := s.Serialize(map[string]interface{}{
			"userID": "fuga",
		})
		req.Header.Set("Authorization", "Bearer "+validJwt)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		auths := strings.SplitN(w.Header().Get("Authorization"), " ", 2)

		claims, err := s.Deserialize(auths[1])
		if err != nil {
			t.Error("deserialize failed:", err)
		}
		if claims["userID"].(string) != "fuga" {
			t.Error("wrong userID filled")
		}
	})

}
