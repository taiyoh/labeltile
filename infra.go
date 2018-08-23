package labeltile

import (
	"github.com/taiyoh/labeltile/app"
	"github.com/taiyoh/labeltile/app/infra"
)

// SetupUserTokenSerializer provides setting UserTokenSerializer into containter
func SetupUserTokenSerializer(container app.Container, conf *Conf) {
	c := conf.JWT
	container.Register("UserTokenSerializer", infra.NewUserTokenSerializer(c.SigningMethod, c.SecretKey, c.ExpireHour))
}

// SetupRepositories provides initializing repositories and setting to container
func SetupRepositories(container app.Container, conf *Conf) {
}
