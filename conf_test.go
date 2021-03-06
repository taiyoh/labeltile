package labeltile_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/taiyoh/labeltile"
)

func TestInvalidConf(t *testing.T) {
	d, _ := os.Getwd()
	d = filepath.Join(d, "test", "conf")

	if _, err := labeltile.NewConf(filepath.Join(d, "none.toml")); err == nil {
		t.Error("file not found")
	}
	if _, err := labeltile.NewConf(filepath.Join(d, "no_jwt_section.toml")); err == nil {
		t.Error("jwt section not found")
	}
	if _, err := labeltile.NewConf(filepath.Join(d, "no_server_section.toml")); err == nil {
		t.Error("server section not found")
	}
	if _, err := labeltile.NewConf(filepath.Join(d, "invalid_jwt_section.toml")); err == nil {
		t.Error("invalid jwt section")
	}
	if _, err := labeltile.NewConf(filepath.Join(d, "invalid_server_section1.toml")); err == nil {
		t.Error("invalid port number")
	}
	if _, err := labeltile.NewConf(filepath.Join(d, "invalid_server_section2.toml")); err == nil {
		t.Error("invalid template path")
	}
	if _, err := labeltile.NewConf(filepath.Join(d, "invalid_oauth2google_section1.toml")); err == nil {
		t.Error("requires client_id and client_secret")
	}
	if _, err := labeltile.NewConf(filepath.Join(d, "invalid_oauth2google_section2.toml")); err == nil {
		t.Error("invalid redirect_url")
	}
}

func TestValidConf(t *testing.T) {
	d, e := ioutil.TempDir("", "labeltile_test")
	if e != nil {
		t.Error("TempDir execution failed:", e)
		return
	}
	defer os.RemoveAll(d)

	f := filepath.Join(d, "valid.toml")
	loadValidConfPath(f, templatePath())
	if _, err := labeltile.NewConf(f); err != nil {
		t.Error("valid section:", err)
	}
}

func templatePath() string {
	currentDir, _ := os.Getwd()
	return filepath.Join(currentDir, "templates")
}

func loadValidConfPath(filePath, templateDir string) {
	valid_conf_text := `
[server]
port=3000
template="{{.template}}"

[jwt]
secret_key="foobar"
signing_method="HS512"
expire_hour=1

[oauth2.google]
client_id="foo"
client_secret="bar"
redirect_url="https://example.com/auth/google/callback"

`

	tmpl := template.Must(template.New("valid_conf").Parse(valid_conf_text))
	bu := &bytes.Buffer{}
	tmpl.Execute(bu, map[string]string{"template": strings.Replace(templateDir, "\\", "\\\\", -1)})
	ioutil.WriteFile(filePath, bu.Bytes(), 0666)
}
