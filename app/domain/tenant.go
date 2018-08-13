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

type tenantMembers []UserID

func (t tenantMembers) Append(m UserID) tenantMembers {
	nt := t[:]
	return append(nt, m)
}

func (t tenantMembers) Delete(m UserID) tenantMembers {
	users := tenantMembers{}
	for _, u := range t {
		if u != m {
			users = append(users, u)
		}
	}
	return users
}

// Tenant manages langs and categories
type Tenant struct {
	ID          TenantID
	Name        string
	DefaultLang LangID
	Members     tenantMembers
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
		Members:     tenantMembers{},
		Languages:   tenantLanguages{dl},
	}
}

// AddLanguage set new language to use
func (t *Tenant) AddLanguage(l LangID) *Tenant {
	return &Tenant{
		ID:          t.ID,
		Name:        t.Name,
		DefaultLang: t.DefaultLang,
		Members:     t.Members,
		Languages:   t.Languages.Append(l),
	}
}

// DeleteLanguage unset language for unuse
func (t *Tenant) DeleteLanguage(l LangID) *Tenant {
	return &Tenant{
		ID:          t.ID,
		Name:        t.Name,
		DefaultLang: t.DefaultLang,
		Members:     t.Members,
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
			Members:     t.Members,
			Languages:   t.Languages,
		}
		return nt, nil
	}
	return nil, errors.New("specified lang is not registered")
}

// AddMember set new member of tenant
func (t *Tenant) AddMember(m UserID) *Tenant {
	return &Tenant{
		ID:          t.ID,
		Name:        t.Name,
		DefaultLang: t.DefaultLang,
		Members:     t.Members.Append(m),
		Languages:   t.Languages,
	}
}

// DeleteMember unset member
func (t *Tenant) DeleteMember(m UserID) *Tenant {
	return &Tenant{
		ID:          t.ID,
		Name:        t.Name,
		DefaultLang: t.DefaultLang,
		Members:     t.Members.Delete(m),
		Languages:   t.Languages,
	}
}

// TenantSpecification provides validation for tenant operation
type TenantSpecification struct {
	tRepo TenantRepository
	uRepo UserRepository
}

// NewTenantSpecification returns TenantSpecification object
func NewTenantSpecification(t TenantRepository, u UserRepository) *TenantSpecification {
	return &TenantSpecification{tRepo: t, uRepo: u}
}

// SpecifyOperateLabel returns whether given operator can operate by given tenant or not
func (s *TenantSpecification) SpecifyOperateLabel(tenantid, opID string) error {
	t := s.tRepo.Find(tenantid)
	if t == nil {
		return errors.New("tenant not found")
	}
	o := s.uRepo.Find(opID)
	if o == nil {
		return errors.New("operator not found")
	}
	isMember := false
	for _, u := range t.Members {
		if u == o.ID {
			isMember = true
			break
		}
	}
	if !isMember {
		return errors.New("operator is not a member of this tenant")
	}

	canEdit := false
	for _, r := range o.Roles {
		if r == RoleEditor {
			canEdit = true
			break
		}
	}
	if !canEdit {
		return errors.New("operator has no permission")
	}

	return nil
}
