package app

import (
	"github.com/taiyoh/labeltile/app/domain"
)

// LabelAddService provides registering label and validation
func LabelAddService(opID, tenantID, key string, urepo domain.UserRepository, lrepo domain.LabelRepository, trepo domain.TenantRepository) error {
	tspec := domain.NewTenantSpecification(trepo, urepo)
	if err := tspec.SpecifyOperateLabel(tenantID, opID); err != nil {
		return err
	}
	tID := domain.TenantID(tenantID)
	lspec := domain.NewLabelSpecification(lrepo)
	if err := lspec.SpecifyAddLabel(tID, key); err != nil {
		return err
	}
	factory := domain.NewLabelFactory(lrepo)
	label := factory.Build(tID, key)
	lrepo.Save(label)
	return nil
}
