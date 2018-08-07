package domain

import (
	"errors"
)

type tenantLanguages []LangID

func (t tenantLanguages) Append(l LangID) tenantLanguages {
	nt := t[:]
	return append(nt, l)
}

func (t tenantLanguages) Delete(l LangID) tenantLanguages {
	langs := tenantLanguages{}
	for _, ln := range t {
		if ln != l {
			langs = append(langs, ln)
		}
	}
	return langs
}

func (t tenantLanguages) Exists(l LangID) bool {
	for _, ln := range t {
		if l == ln {
			return true
		}
	}
	return false
}

type tenantCategories []CategoryID

func (t tenantCategories) Append(c CategoryID) tenantCategories {
	nt := t[:]
	return append(nt, c)
}

func (t tenantCategories) Delete(c CategoryID) tenantCategories {
	cats := tenantCategories{}
	for _, ca := range t {
		if ca != c {
			cats = append(cats, ca)
		}
	}
	return cats
}

// Tenant manages langs and categories
type Tenant struct {
	ID          TenantID
	Name        string
	DefaultLang LangID
	Languages   tenantLanguages
	Categories  tenantCategories
}

// NewTenant returns initialized Tenant object
func NewTenant(id TenantID, name string, dl LangID) *Tenant {
	return &Tenant{
		ID:          id,
		Name:        name,
		DefaultLang: dl,
		Languages:   tenantLanguages{dl},
		Categories:  tenantCategories{},
	}
}

// AddLanguage set new language to use
func (t *Tenant) AddLanguage(l LangID) *Tenant {
	return &Tenant{
		ID:          t.ID,
		Name:        t.Name,
		DefaultLang: t.DefaultLang,
		Languages:   t.Languages.Append(l),
		Categories:  t.Categories,
	}
}

// DeleteLanguage unset language for unuse
func (t *Tenant) DeleteLanguage(l LangID) *Tenant {
	return &Tenant{
		ID:          t.ID,
		Name:        t.Name,
		DefaultLang: t.DefaultLang,
		Languages:   t.Languages.Delete(l),
		Categories:  t.Categories,
	}
}

// ChangeDefaultLang set default language for tenant
func (t *Tenant) ChangeDefaultLang(l LangID) (*Tenant, error) {
	if t.Languages.Exists(l) {
		nt := &Tenant{
			ID:          t.ID,
			Name:        t.Name,
			DefaultLang: l,
			Languages:   t.Languages,
			Categories:  t.Categories,
		}
		return nt, nil
	}
	return nil, errors.New("specified lang is not registered")
}

// AddCategory set new category for managing in tenant
func (t *Tenant) AddCategory(id CategoryID) *Tenant {
	return &Tenant{
		ID:          t.ID,
		Name:        t.Name,
		DefaultLang: t.DefaultLang,
		Languages:   t.Languages,
		Categories:  t.Categories.Append(id),
	}
}

// DeleteCategory unset category
func (t *Tenant) DeleteCategory(id CategoryID) *Tenant {
	return &Tenant{
		ID:          t.ID,
		Name:        t.Name,
		DefaultLang: t.DefaultLang,
		Languages:   t.Languages,
		Categories:  t.Categories.Delete(id),
	}
}
