package domain_test

import (
	"testing"

	"github.com/taiyoh/labeltile/app/domain"
)

func TestTenant(t *testing.T) {
	tenant := domain.NewTenant(domain.TenantID("1"), "foo", domain.LangID("ja"))
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

	tenant = tenant.AddCategory(domain.CategoryID("p1"))
	tenant = tenant.AddCategory(domain.CategoryID("p2"))
	if len(tenant.Categories) != 2 {
		t.Error("category:p1 and category:p2 should be added")
	}
	if tenant.Categories[0] != domain.CategoryID("p1") {
		t.Error("category:p1 should be added")
	}
	if tenant.Categories[1] != domain.CategoryID("p2") {
		t.Error("category:p2 should be added")
	}
	tenant = tenant.DeleteCategory(domain.CategoryID("p1"))
	if len(tenant.Categories) != 1 {
		t.Error("category:p1 should be removed")
	}
	if tenant.Categories[0] != domain.CategoryID("p2") {
		t.Error("category:p2 should be set")
	}
}
