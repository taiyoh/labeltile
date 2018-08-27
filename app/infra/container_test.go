package infra_test

import (
	"testing"

	"github.com/taiyoh/labeltile/app/infra"
)

func TestNewContainer(t *testing.T) {
	c := infra.NewContainer()
	if roleRepo := c.RoleRepository(); roleRepo == nil {
		t.Error("Role Repository not returns")
	}

	if s := c.UserTokenSerializer(); s != nil {
		t.Error("UserTokenSerializer returns")
	}

	c.Register("UserTokenSerializer", infra.NewUserTokenSerializer("HS512", "foobar", 1))
	if s := c.UserTokenSerializer(); s == nil {
		t.Error("UserTokenSerializer not returns")
	}

}
