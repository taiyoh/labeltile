package mock

import "github.com/taiyoh/labeltile/app/domain"

type LabelRepositoryImpl struct {
	domain.LabelRepository
	DispenseIDFunc func() domain.LabelID
	Labels         map[domain.LabelID]*domain.Label
}

func LoadLabelRepoImpl(f func() domain.LabelID) *LabelRepositoryImpl {
	return &LabelRepositoryImpl{DispenseIDFunc: f, Labels: map[domain.LabelID]*domain.Label{}}
}

func (r *LabelRepositoryImpl) DispenseID() domain.LabelID {
	return r.DispenseIDFunc()
}

func (r *LabelRepositoryImpl) Find(id string) *domain.Label {
	l, _ := r.Labels[domain.LabelID(id)]
	return l
}

func (r *LabelRepositoryImpl) FindByKey(key string, tenantID domain.TenantID) *domain.Label {
	for _, label := range r.Labels {
		if label.Tenant == tenantID && label.Key == key {
			return label
		}
	}
	return nil
}

func (r *LabelRepositoryImpl) Save(l *domain.Label) {
	r.Labels[l.ID] = l
}
