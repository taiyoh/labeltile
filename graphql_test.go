package labeltile_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/taiyoh/labeltile"
)

func newReader(s string) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewBufferString(s))
}

func TestBrokenRequest(t *testing.T) {
	if _, err := labeltile.NewGraphQLRequest(newReader(`{"foo":[}`)); err == nil {
		t.Error("broken request")
	}

	if _, err := labeltile.NewGraphQLRequest(newReader(`{"foo":"bar"}`)); err == nil {
		t.Error("requires query and variables")
	}
}

func TestNewRequest(t *testing.T) {
	reqStr := `{"variables": {}, "query": "query { operator { id } }"}`
	req, err := labeltile.NewGraphQLRequest(newReader(reqStr))
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
