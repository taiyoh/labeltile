package app_test

import (
	"testing"

	"github.com/taiyoh/labeltile/app"

	"github.com/taiyoh/labeltile/app/domain"
	"github.com/taiyoh/labeltile/app/infra/mock"
)

func TestLabelAddService(t *testing.T) {
	uid := "1"
	urepo := mock.LoadUserRepoImpl(func() domain.UserID {
		return domain.UserID(uid)
	})
	ufactory := domain.NewUserFactory(urepo)
	user := ufactory.Build(domain.UserMail("foo@example.com"))
	user = user.AddRole(domain.RoleEditor)
	urepo.Save(user)

	tid := "1"
	trepo := mock.LoadTenantRepoImpl(func() domain.TenantID {
		return domain.TenantID(tid)
	})
	tfactory := domain.NewTenantFactory(trepo)
	tenant := tfactory.Build("foo", domain.LangID("ja"))
	tenant = tenant.AddMember(user.ID)
	trepo.Save(tenant)

	lid := "1"
	lrepo := mock.LoadLabelRepoImpl(func() domain.LabelID {
		return domain.LabelID(lid)
	})

	if err := app.LabelAddService(uid, "2", "foo", urepo, lrepo, trepo); err == nil {
		t.Error("tenant not found")
	}
	if err := app.LabelAddService(uid, tid, "foo", urepo, lrepo, trepo); err != nil {
		t.Error("operation is valid")
	}
	if err := app.LabelAddService(uid, tid, "foo", urepo, lrepo, trepo); err == nil {
		t.Error("label:foo already registered")
	}

}
