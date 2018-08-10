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

// Tenant manages langs and categories
type Tenant struct {
	ID          TenantID
	Name        string
	DefaultLang LangID
	Languages   tenantLanguages
}

// TenantRepository is interface for Tenant repository
type TenantRepository interface {
	DispenseID() TenantID
	Find(id string) *Tenant
	Save(t *Tenant)
}

// TenantFactory provides builder for Tenant
type TenantFactory struct {
	tRepo TenantRepository
}

// NewTenantFactory returns TenantFactory object
func NewTenantFactory(r TenantRepository) *TenantFactory {
	return &TenantFactory{tRepo: r}
}

// Build returns initialized Tenant object
func (f *TenantFactory) Build(name string, dl LangID) *Tenant {
	return &Tenant{
		ID:          f.tRepo.DispenseID(),
		Name:        name,
		DefaultLang: dl,
		Languages:   tenantLanguages{dl},
	}
}

// AddLanguage set new language to use
func (t *Tenant) AddLanguage(l LangID) *Tenant {
	return &Tenant{
		ID:          t.ID,
		Name:        t.Name,
		DefaultLang: t.DefaultLang,
		Languages:   t.Languages.Append(l),
	}
}

// DeleteLanguage unset language for unuse
func (t *Tenant) DeleteLanguage(l LangID) *Tenant {
	return &Tenant{
		ID:          t.ID,
		Name:        t.Name,
		DefaultLang: t.DefaultLang,
		Languages:   t.Languages.Delete(l),
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
		}
		return nt, nil
	}
	return nil, errors.New("specified lang is not registered")
}
