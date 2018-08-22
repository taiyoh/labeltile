package mock

import (
	"github.com/taiyoh/labeltile/app/domain"
	"github.com/taiyoh/labeltile/app/infra"
)

func LoadContainer() *infra.Container {
	labelID := 1
	labelRepo := LoadLabelRepoImpl(func() domain.LabelID {
		lid := labelID
		labelID++
		return domain.LabelID(string(lid))
	})
	userID := 1
	userRepo := LoadUserRepoImpl(func() domain.UserID {
		uid := userID
		userID++
		return domain.UserID(string(uid))
	})
	tenantID := 1
	tenantRepo := LoadTenantRepoImpl(func() domain.TenantID {
		tid := tenantID
		tenantID++
		return domain.TenantID(string(tid))
	})

	c := infra.NewContainer()
	c.Register("LabelRepository", labelRepo)
	c.Register("UserRepository", userRepo)
	c.Register("TenantRepository", tenantRepo)

	return c
}
