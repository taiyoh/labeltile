package labeltile_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/taiyoh/labeltile"
)

func TestNewLabeltile(t *testing.T) {

	d, _ := ioutil.TempDir("", "labeltile_test")
	defer os.RemoveAll(d)

	f := filepath.Join(d, "valid.toml")
	loadValidConfPath(f, templatePath())
	conf, _ := labeltile.NewConf(f)

	l := labeltile.NewLabeltile(
		conf,
		labeltile.SetupUserTokenSerializer,
		labeltile.SetupRepositories,
	)
	l.SetupRoutes(
		labeltile.SetupHomeHandler,
		labeltile.SetupOAuth2GoogleHandler,
		labeltile.SetupGraphQLHandler,
	)
}
