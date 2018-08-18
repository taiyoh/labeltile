package infra

import (
	"github.com/taiyoh/labeltile/app/domain"
)

type Container struct {
	serializer       *UserTokenSerializer
	userRepository   domain.UserRepository
	labelRepository  domain.LabelRepository
	tenantRepository domain.TenantRepository
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) SetUserTokenSerializer(s *UserTokenSerializer) {
	c.serializer = s
}

func (c *Container) UserTokenSerializer() *UserTokenSerializer {
	return c.serializer
}
