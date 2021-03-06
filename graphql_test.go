package labeltile_test

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"testing"

	"github.com/taiyoh/labeltile"
	"github.com/taiyoh/labeltile/app"
	"github.com/taiyoh/labeltile/app/domain"
	"github.com/taiyoh/labeltile/app/infra/mock"
)

func newReader(s string) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewBufferString(s))
}

func TestBrokenRequest(t *testing.T) {
	s := labeltile.InitializeGraphQLSchema(mock.LoadContainer())
	if _, err := labeltile.NewGraphQLRequest(s, newReader(`{"foo":[}`)); err == nil {
		t.Error("broken request")
	}

	if _, err := labeltile.NewGraphQLRequest(s, newReader(`{"foo":"bar"}`)); err == nil {
		t.Error("requires query and variables")
	}
}

func TestNewRequest(t *testing.T) {
	s := labeltile.InitializeGraphQLSchema(mock.LoadContainer())
	reqStr := `{"variables": {}, "query": "query { operator { id } }"}`
	req, err := labeltile.NewGraphQLRequest(s, newReader(reqStr))
	if err != nil {
		t.Error("error found: " + err.Error())
	}
	if req.Query != "query { operator { id } }" {
		t.Error("query is wrong")
	}
	if len(req.Variables) != 0 {
		t.Error("something variable is in")
	}
}

func TestRunGraphQLRequest(t *testing.T) {
	c := mock.LoadContainer()
	s := labeltile.InitializeGraphQLSchema(c)
	reqStr := `{"variables": {}, "query": "query { operator { id mail } }"}`
	req, _ := labeltile.NewGraphQLRequest(s, newReader(reqStr))

	factory := domain.NewUserFactory(c.UserRepository())
	u := factory.Build(domain.UserMail("foo@example.com"))
	u = u.AddRole(domain.RoleEditor)
	c.UserRepository().Save(u)

	ctx := context.Background()

	t.Run("not loggedin", func(t *testing.T) {
		res := req.Run(ctx)
		data := res["data"].(map[string]interface{})
		op, ok := data["operator"]
		if !ok {
			t.Error("operator key not found")
		}
		if op != nil {
			t.Error("operator data found")
		}
	})

	t.Run("already loggedin", func(t *testing.T) {
		res := req.Run(context.WithValue(ctx, app.UserIDCtxKey, string(u.ID)))
		data := res["data"].(map[string]interface{})
		op, ok := data["operator"].(map[string]interface{})
		if !ok {
			t.Error("operator key not found")
		}
		if v, exists := op["id"]; !exists || v.(string) != string(u.ID) {
			t.Error("wrong id fetched:", v)
		}
		if v, exists := op["mail"]; !exists || v.(string) != string(u.Mail) {
			t.Error("wrong mail fetched:", v)
		}
		if _, exists := op["roles"]; exists {
			t.Error("unspecified field returns")
		}
	})

}
