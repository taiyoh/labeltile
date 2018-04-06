package domain

import (
	"errors"
)

// Tenant manages langs and categories
type Tenant struct {
	ID          TenantID
	Name        string
	DefaultLang LangID
	Languages   []LangID
	Categories  []CategoryID
}

// NewTenant returns initialized Tenant object
func NewTenant(id TenantID, name string, dl LangID) *Tenant {
	return &Tenant{
		ID:          id,
		Name:        name,
		DefaultLang: dl,
		Languages:   []LangID{dl},
	}
}

// AddLanguage set new language to use
func (t *Tenant) AddLanguage(l LangID) {
	t.Languages = append(t.Languages, l)
}

// DeleteLanguage unset language for unuse
func (t *Tenant) DeleteLanguage(l LangID) {
	langs := []LangID{}
	for _, ln := range t.Languages {
		if ln != l {
			langs = append(langs, ln)
		}
	}
	t.Languages = langs
}

// ChangeDefaultLang set default language for tenant
func (t *Tenant) ChangeDefaultLang(l LangID) error {
	for _, ln := range t.Languages {
		if ln == l {
			t.DefaultLang = l
			return nil
		}
	}
	return errors.New("specified lang is not registered")
}

// AddCategory set new category for managing in tenant
func (t *Tenant) AddCategory(id CategoryID) {
	t.Categories = append(t.Categories, id)
}

// DeleteCategory unset category
func (t *Tenant) DeleteCategory(id CategoryID) {
	cats := []CategoryID{}
	for _, ca := range t.Categories {
		if ca != id {
			cats = append(cats, ca)
		}
	}
	t.Categories = cats
}
