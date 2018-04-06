package domain

import "testing"

func TestTenant(t *testing.T) {
	tenant := NewTenant(TenantID("1"), "foo", LangID("ja"))
	if tenant.Name != "foo" {
		t.Error("tenant name should be 'foo'")
	}
	if tenant.DefaultLang != LangID("ja") {
		t.Error("tenant default lang is not set")
	}
	if len(tenant.Languages) != 1 {
		t.Error("tenant langs should be default lang only")
	}
	if tenant.Languages[0] != LangID("ja") {
		t.Error("tenant langs is wrong")
	}

	tenant.AddLanguage(LangID("en"))
	if len(tenant.Languages) != 2 {
		t.Error("lang:en should be added")
	}
	if tenant.Languages[1] != LangID("en") {
		t.Error("lang:en should be added")
	}

	if err := tenant.ChangeDefaultLang(LangID("fr")); err == nil {
		t.Error("lang:fr is not set")
	}
	if err := tenant.ChangeDefaultLang(LangID("en")); err != nil {
		t.Error("lang:en should be set to default lang")
	}
	if tenant.DefaultLang != LangID("en") {
		t.Error("lang:en should be set to default lang")
	}
	tenant.DeleteLanguage(LangID("ja"))
	if len(tenant.Languages) != 1 {
		t.Error("lang:ja should be removed")
	}
	if tenant.Languages[0] != LangID("en") {
		t.Error("lang:ja should be removed")
	}

	tenant.AddCategory(CategoryID("p1"))
	tenant.AddCategory(CategoryID("p2"))
	if len(tenant.Categories) != 2 {
		t.Error("category:p1 and category:p2 should be added")
	}
	if tenant.Categories[0] != CategoryID("p1") {
		t.Error("category:p1 should be added")
	}
	if tenant.Categories[1] != CategoryID("p2") {
		t.Error("category:p2 should be added")
	}
	tenant.DeleteCategory(CategoryID("p1"))
	if len(tenant.Categories) != 1 {
		t.Error("category:p1 should be removed")
	}
	if tenant.Categories[0] != CategoryID("p2") {
		t.Error("category:p2 should be set")
	}
}
