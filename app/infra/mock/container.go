package mock

import (
	"github.com/taiyoh/labeltile/app"
	"github.com/taiyoh/labeltile/app/domain"
)

type ContainerImpl struct {
	app.Container
	labelRepo  domain.LabelRepository
	userRepo   domain.UserRepository
	tenantRepo domain.TenantRepository
	serializer app.UserTokenSerializer
	roleRepo   *domain.RoleRepository
}

func LoadContainerImpl() *ContainerImpl {
	labelID := 1
	userID := 1
	tenantID := 1
	return &ContainerImpl{
		labelRepo: LoadLabelRepoImpl(func() domain.LabelID {
			lid := labelID
			labelID++
			return domain.LabelID(string(lid))
		}),
		userRepo: LoadUserRepoImpl(func() domain.UserID {
			uid := userID
			userID++
			return domain.UserID(string(uid))
		}),
		tenantRepo: LoadTenantRepoImpl(func() domain.TenantID {
			tid := tenantID
			tenantID++
			return domain.TenantID(string(tid))
		}),
		roleRepo: &domain.RoleRepository{},
	}
}

func (c *ContainerImpl) SetSerializer(s app.UserTokenSerializer) {
	c.serializer = s
}

func (c *ContainerImpl) UserTokenSerializer() app.UserTokenSerializer {
	return c.serializer
}

func (c *ContainerImpl) RoleRepository() *domain.RoleRepository {
	return c.roleRepo
}

func (c *ContainerImpl) UserRepository() domain.UserRepository {
	return c.userRepo
}
