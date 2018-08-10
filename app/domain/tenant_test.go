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
