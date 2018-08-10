package mock

import "github.com/taiyoh/labeltile/app/domain"

type TenantRepositoryImpl struct {
	domain.TenantRepository
	DispenseIDFunc func() domain.TenantID
	Tenants        map[domain.TenantID]*domain.Tenant
}

func LoadTenantRepoImpl(f func() domain.TenantID) *TenantRepositoryImpl {
	return &TenantRepositoryImpl{DispenseIDFunc: f, Tenants: map[domain.TenantID]*domain.Tenant{}}
}

func (r *TenantRepositoryImpl) DispenseID() domain.TenantID {
	return r.DispenseIDFunc()
}

func (r *TenantRepositoryImpl) Find(id string) *domain.Tenant {
	t, _ := r.Tenants[domain.TenantID(id)]
	return t
}

func (r *TenantRepositoryImpl) Save(t *domain.Tenant) {
	r.Tenants[t.ID] = t
}
