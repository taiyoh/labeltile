package infra

import (
	"github.com/taiyoh/labeltile/app"
	"github.com/taiyoh/labeltile/app/domain"
)

type Container struct {
	app.Container
	serializer       app.UserTokenSerializer
	oauth2Google     app.OAuth2Google
	userRepository   domain.UserRepository
	labelRepository  domain.LabelRepository
	tenantRepository domain.TenantRepository
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) SetUserTokenSerializer(s app.UserTokenSerializer) {
	c.serializer = s
}

func (c *Container) UserTokenSerializer() app.UserTokenSerializer {
	return c.serializer
}

func (c *Container) OAuth2Google() app.OAuth2Google {
	return c.oauth2Google
}
