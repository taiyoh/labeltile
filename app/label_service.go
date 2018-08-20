package app

import (
	"github.com/taiyoh/labeltile/app/domain"
)

// LabelAddService provides registering label and validation
func LabelAddService(opID, tenantID, key string, container interface {
	UserRepository() domain.UserRepository
	LabelRepository() domain.LabelRepository
	TenantRepository() domain.TenantRepository
}) error {
	tspec := domain.NewTenantSpecification(container.TenantRepository(), container.UserRepository())
	if err := tspec.SpecifyOperateLabel(tenantID, opID); err != nil {
		return err
	}
	tID := domain.TenantID(tenantID)
	lspec := domain.NewLabelSpecification(container.LabelRepository())
	if err := lspec.SpecifyAddLabel(tID, key); err != nil {
		return err
	}
	factory := domain.NewLabelFactory(container.LabelRepository())
	label := factory.Build(tID, key)
	container.LabelRepository().Save(label)
	return nil
}
