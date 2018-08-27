package labeltile_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/taiyoh/labeltile/app/domain"

	"github.com/gin-gonic/gin"
	"github.com/taiyoh/labeltile"
	"github.com/taiyoh/labeltile/app/infra"
	"github.com/taiyoh/labeltile/app/infra/mock"
)

func TestHomeHandler(t *testing.T) {
	c := mock.LoadContainer()
	r := gin.Default()
	d, _ := os.Getwd()
	tmplPath := filepath.Join(d, "templates", "index.tmpl")
	r.LoadHTMLFiles(tmplPath)
	labeltile.SetupHomeHandler(r, c)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Error("response failed")
	}
	buf := bytes.NewBuffer([]byte{})
	buf.ReadFrom(res.Body)

	tmplContent, _ := ioutil.ReadFile(tmplPath)
	if bytes.Compare(buf.Bytes(), tmplContent) != 0 {
		t.Error("wrong content returns")
	}
}

func TestGraphQLRequest(t *testing.T) {
	c := mock.LoadContainer()
	c.Register("UserTokenSerializer", infra.NewUserTokenSerializer("HS512", "foobar", 1))
	r := gin.Default()
	r.Use(labeltile.SetupUserTokenMiddleware(c))
	labeltile.SetupGraphQLHandler(r, c)

	c.Register("UserTokenSerializer", infra.NewUserTokenSerializer("HS512", "foobar", 1))

	factory := domain.NewUserFactory(c.UserRepository())
	user := factory.Build(domain.UserMail("foo@example.com"))
	c.UserRepository().Save(user)

	t.Run("wrong Content-Type", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/graphql", nil)
		r.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()
		if res.StatusCode != http.StatusBadRequest {
			t.Error("unexpected response code")
		}
		if !strings.HasPrefix(res.Header.Get("Content-Type"), "application/json") {
			t.Error("response is not json:", res.Header.Get("Content-Type"))
		}
		resBuf := bytes.NewBuffer([]byte{})
		resBuf.ReadFrom(res.Body)
		resStruct := &struct {
			Error string
		}{}
		json.Unmarshal(resBuf.Bytes(), resStruct)
		if resStruct.Error != "bad Content-Type" {
			t.Error("wrong error returns")
		}
	})

	t.Run("broken request", func(t *testing.T) {
		w := httptest.NewRecorder()
		buf := bytes.NewBufferString(`{"variables":{}:,,,"query":"}`)
		req, _ := http.NewRequest("POST", "/graphql", buf)
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()
		if res.StatusCode != http.StatusBadRequest {
			t.Error("unexpected response code")
		}
		resBuf := bytes.NewBuffer([]byte{})
		resBuf.ReadFrom(res.Body)
		resStruct := &struct {
			Error string
		}{}
		json.Unmarshal(resBuf.Bytes(), resStruct)
		if resStruct.Error != "broken request" {
			t.Error("wrong error returns")
		}
	})

	t.Run("right request", func(t *testing.T) {
		w := httptest.NewRecorder()
		buf := bytes.NewBufferString(`{"variables":{},"query":"query { operator { id } }"}`)
		req, _ := http.NewRequest("POST", "/graphql", buf)
		req.Header.Set("Content-Type", "application/json")
		cl := c.UserTokenSerializer().NewClaims()
		cl.UserID(string(user.ID))
		token, _ := c.UserTokenSerializer().Serialize(cl)
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			t.Error("unexpected response code:", res.StatusCode)
		}
	})

}
