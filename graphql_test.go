package labeltile_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/taiyoh/labeltile"
	"github.com/taiyoh/labeltile/app/infra"
)

func loadSerializer() *infra.UserTokenSerializer {
	return infra.NewUserTokenSerializer("HS512", "foobar", 1)
}

func newReader(s string) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewBufferString(s))
}

func TestBrokenRequest(t *testing.T) {
	c := infra.NewContainer()
	c.SetUserTokenSerializer(loadSerializer())
	if _, err := labeltile.NewGraphQLRequest(newReader(`{"foo":[}`), "", c); err == nil {
		t.Error("broken request")
	}

	if _, err := labeltile.NewGraphQLRequest(newReader(`{"foo":"bar"}`), "", c); err == nil {
		t.Error("requires query and variables")
	}
}

func TestNewRequestWithoutToken(t *testing.T) {
	c := infra.NewContainer()
	c.SetUserTokenSerializer(loadSerializer())
	reqStr := `{"variables": {}, "query": "query { operator { id } }"}`
	req, err := labeltile.NewGraphQLRequest(newReader(reqStr), "", c)
	if err != nil {
		t.Error("error found: " + err.Error())
	}
	if req.User != nil {
		t.Error("user should not be exists")
	}
	if req.Query != "query { operator { id } }" {
		t.Error("query is wrong")
	}
	if len(req.Variables) != 0 {
		t.Error("something variable is in")
	}
}

func TestNewRequestWithToken(t *testing.T) {
	c := infra.NewContainer()
	c.SetUserTokenSerializer(loadSerializer())
	reqStr := `{"variables": {}, "query": "query { operator { id } }"}`
	_, err := labeltile.NewGraphQLRequest(newReader(reqStr), "hoge", c)
	if err == nil {
		t.Error("user token is wrong")
	}
	token, _ := c.UserTokenSerializer().Serialize(map[string]interface{}{"userID": "nya-"})
	req, err := labeltile.NewGraphQLRequest(newReader(reqStr), token, c)
	if err != nil || req.User == nil {
		t.Error("user token is valid")
	}
	if req.User.ID != "nya-" {
		t.Error("userID is wrong")
	}
}
