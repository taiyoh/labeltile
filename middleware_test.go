package labeltile_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/taiyoh/labeltile"
	"github.com/taiyoh/labeltile/app"
	"github.com/taiyoh/labeltile/app/infra"
	"github.com/taiyoh/labeltile/app/infra/mock"

	"github.com/gin-gonic/gin"
)

func TestTokenMiddleware(t *testing.T) {
	router := gin.Default()
	c := mock.LoadContainer()
	s := infra.NewUserTokenSerializer("HS512", "foobar", 1)
	c.Register("UserTokenSerializer", s)
	router.Use(labeltile.SetupUserTokenMiddleware(c))
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
		if w.Header().Get("Authorization") != "" {
			t.Error("unknown Authorization header returns")
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
		id := claims.FindUserID()
		if id != "foo" {
			t.Error("wrong userID")
		}
	})
	t.Run("broken Authorization header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		brokenJwt := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		req.Header.Set("Authorization", "Bearer "+brokenJwt)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Header().Get("Authorization") != "" {
			t.Error("unknown Authorization header returns")
		}
	})
	t.Run("valid Authorization header", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		cl := s.NewClaims()
		cl.UserID("fuga")
		validJwt, _ := s.Serialize(cl)
		req.Header.Set("Authorization", "Bearer "+validJwt)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		auths := strings.SplitN(w.Header().Get("Authorization"), " ", 2)

		claims, err := s.Deserialize(auths[1])
		if err != nil {
			t.Error("deserialize failed:", err)
		}
		if claims.FindUserID() != cl.FindUserID() {
			t.Error("wrong userID filled")
		}
	})

}

func TestSessionMiddleware(t *testing.T) {
	router := gin.Default()
	c := mock.LoadContainer()
	s := mock.LoadSession("foo", 3600)
	c.Register("SessionStorage", s)

	sessionID := ""

	router.Use(
		func(c *gin.Context) {
			if sessionID != "" {
				c.Set("sessionID", sessionID)
			}
			c.Next()
		},
		labeltile.SetupSessionMiddleware(c),
	)
	router.GET("/", func(c *gin.Context) {
		if se, exists := c.Get("session"); exists {
			session, _ := se.(app.SessionData)
			session.Set("foo", []int{4423, 6061})
		}
		c.JSON(http.StatusOK, gin.H{"sessionID": sessionID})
	})

	t.Run("no session id assigned", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		res := &struct {
			SessionID string `json:"sessionID"`
		}{}
		json.Unmarshal(w.Body.Bytes(), res)
		if res.SessionID != "" {
			t.Error("wrong session id assigned")
		}
	})
	sessionID = "foobar"
	t.Run("session id assigned", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		res := &struct {
			SessionID string `json:"sessionID"`
		}{}
		json.Unmarshal(w.Body.Bytes(), res)
		if res.SessionID != "foobar" {
			t.Error("wrong session id assigned")
		}
		sessionData := s.Find("foobar")
		if sessionData == nil {
			t.Error("session not found")
			return
		}
		d := sessionData.Get("foo")
		if d == nil {
			t.Error("session data not found")
			return
		}
		ints, ok := d.([]int)
		if !ok {
			t.Error("stored 'foo' is wrong")
			return
		}
		if len(ints) != 2 {
			t.Error("contianed []int is wrong")
		}
		if ints[0] != 4423 || ints[1] != 6061 {
			t.Error("contained []int is wrong")
		}
	})
}
