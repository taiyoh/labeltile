package infra

import (
	"github.com/taiyoh/labeltile/app"
	"github.com/taiyoh/labeltile/app/domain"
)

// Container is object stores using Service Locator pattern
type Container struct {
	app.Container
	stores map[string]interface{}
}

// NewContainer returns Container object
func NewContainer() *Container {
	c := &Container{stores: map[string]interface{}{}}
	c.Register("RoleRepository", &domain.RoleRepository{})
	return c
}

// Register provides store given object by given name
func (c *Container) Register(name string, obj interface{}) {
	c.stores[name] = obj
}

// UserTokenSerializer is interface for fetching app.UserTokenSerializer from container stores
func (c *Container) UserTokenSerializer() app.UserTokenSerializer {
	o, ook := c.stores["UserTokenSerializer"].(app.UserTokenSerializer)
	if !ook {
		return nil
	}
	return o
}

// OAuth2Google is interface for fetching app.OAuth2Google from container stores
func (c *Container) OAuth2Google() app.OAuth2Google {
	o, ook := c.stores["OAuth2Google"].(app.OAuth2Google)
	if !ook {
		return nil
	}
	return o
}

// UserRepository is interface for fetching domain.UserRepository from container stores
func (c *Container) UserRepository() domain.UserRepository {
	o, ook := c.stores["UserRepository"].(domain.UserRepository)
	if !ook {
		return nil
	}
	return o
}

// RoleRepository is interface for fetching domain.RoleRepository from container stores
func (c *Container) RoleRepository() *domain.RoleRepository {
	o, _ := c.stores["RoleRepository"].(*domain.RoleRepository)
	return o
}

// TenantRepository is interface for fetching domain.TenantRepository from container stores
func (c *Container) TenantRepository() domain.TenantRepository {
	o, ook := c.stores["TenantRepository"].(domain.TenantRepository)
	if !ook {
		return nil
	}
	return o
}

// LabelRepository is interface for fetching domain.LabelRepository from container stores
func (c *Container) LabelRepository() domain.LabelRepository {
	o, ook := c.stores["LabelRepository"].(domain.LabelRepository)
	if !ook {
		return nil
	}
	return o
}
