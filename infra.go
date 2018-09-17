package labeltile

import (
	"github.com/taiyoh/labeltile/app"
	"github.com/taiyoh/labeltile/app/infra"
	"github.com/taiyoh/labeltile/app/infra/database"
)

// SetupUserTokenSerializer provides setting UserTokenSerializer into containter
func SetupUserTokenSerializer(container app.Container, conf *Conf) {
	c := conf.JWT
	container.Register("UserTokenSerializer", infra.NewUserTokenSerializer(c.SigningMethod, c.SecretKey, c.ExpireHour))
}

// SetupRepositories provides initializing repositories and setting to container
func SetupRepositories(container app.Container, conf *Conf) {
	c := conf.Database
	db, err := database.New(c.Driver, c.Dsn)
	if err != nil {
		panic(err)
	}
	container.Register("UserRepository", infra.NewUserRepository(db))
}
