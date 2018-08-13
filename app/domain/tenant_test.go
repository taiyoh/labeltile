package domain_test

import (
	"testing"

	"github.com/taiyoh/labeltile/app/domain"
	"github.com/taiyoh/labeltile/app/infra/mock"
)

func TestTenant(t *testing.T) {
	trepo := mock.LoadTenantRepoImpl(func() domain.TenantID {
		return domain.TenantID("1")
	})
	factory := domain.NewTenantFactory(trepo)
	tenant := factory.Build("foo", domain.LangID("ja"))
	if tenant.Name != "foo" {
		t.Error("tenant name should be 'foo'")
	}
	if tenant.DefaultLang != domain.LangID("ja") {
		t.Error("tenant default lang is not set")
	}
	if len(tenant.Languages) != 1 {
		t.Error("tenant langs should be default lang only")
	}
	if tenant.Languages[0] != domain.LangID("ja") {
		t.Error("tenant langs is wrong")
	}

	tenant = tenant.AddLanguage(domain.LangID("en"))
	if len(tenant.Languages) != 2 {
		t.Error("lang:en should be added")
	}
	if tenant.Languages[1] != domain.LangID("en") {
		t.Error("lang:en should be added")
	}

	if _, err := tenant.ChangeDefaultLang(domain.LangID("fr")); err == nil {
		t.Error("lang:fr is not set")
	}
	if tenant, _ = tenant.ChangeDefaultLang(domain.LangID("en")); tenant == nil {
		t.Error("lang:en should be set to default lang")
	}
	if tenant.DefaultLang != domain.LangID("en") {
		t.Error("lang:en should be set to default lang")
	}
	tenant = tenant.DeleteLanguage(domain.LangID("ja"))
	if len(tenant.Languages) != 1 {
		t.Error("lang:ja should be removed")
	}
	if tenant.Languages[0] != domain.LangID("en") {
		t.Error("lang:ja should be removed")
	}
}

func TestTenantSpecification(t *testing.T) {
	tid := "1"
	trepo := mock.LoadTenantRepoImpl(func() domain.TenantID {
		return domain.TenantID(tid)
	})
	tfactory := domain.NewTenantFactory(trepo)
	tenant := tfactory.Build("foo", domain.LangID("ja"))
	trepo.Save(tenant)
	uid := "1"
	urepo := mock.LoadUserRepoImpl(func() domain.UserID {
		return domain.UserID(uid)
	})
	ufactory := domain.NewUserFactory(urepo)
	user := ufactory.Build(domain.UserMail("foo@example.com"))
	urepo.Save(user)

	spec := domain.NewTenantSpecification(trepo, urepo)
	if err := spec.SpecifyOperateLabel("2", uid); err == nil {
		t.Error("tenant not found")
	}
	if err := spec.SpecifyOperateLabel(tid, "2"); err == nil {
		t.Error("operator not found")
	}
	if err := spec.SpecifyOperateLabel(tid, uid); err == nil {
		t.Error("operator is not a member of given tenant")
	}

	tenant = tenant.AddMember(user.ID)
	trepo.Save(tenant)

	if err := spec.SpecifyOperateLabel(tid, uid); err == nil {
		t.Error("operator has no permission for edit")
	}

	user = user.AddRole(domain.RoleEditor)
	urepo.Save(user)

	if err := spec.SpecifyOperateLabel(tid, uid); err != nil {
		t.Error("this operation should be valid")
	}

	tenant = tenant.DeleteMember(user.ID)
	trepo.Save(tenant)

	if err := spec.SpecifyOperateLabel(tid, uid); err == nil {
		t.Error("operator has no permission for edit")
	}

}
