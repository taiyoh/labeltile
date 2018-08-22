package app_test

import (
	"testing"

	"github.com/taiyoh/labeltile/app"

	"github.com/taiyoh/labeltile/app/domain"
	"github.com/taiyoh/labeltile/app/infra/mock"
)

func TestLabelAddService(t *testing.T) {
	c := mock.LoadContainer()
	ufactory := domain.NewUserFactory(c.UserRepository())
	user := ufactory.Build(domain.UserMail("foo@example.com"))
	user = user.AddRole(domain.RoleEditor)
	c.UserRepository().Save(user)
	uid := string(user.ID)

	tfactory := domain.NewTenantFactory(c.TenantRepository())
	tenant := tfactory.Build("foo", domain.LangID("ja"))
	tenant = tenant.AddMember(user.ID)
	c.TenantRepository().Save(tenant)
	tid := string(tenant.ID)

	if err := app.LabelAddService(uid, "2", "foo", c); err == nil {
		t.Error("tenant not found")
	}
	if err := app.LabelAddService(uid, tid, "foo", c); err != nil {
		t.Error("operation is valid")
	}
	if err := app.LabelAddService(uid, tid, "foo", c); err == nil {
		t.Error("label:foo already registered")
	}

}
