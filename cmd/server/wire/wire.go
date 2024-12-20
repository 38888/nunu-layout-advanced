//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/38888/nunu-layout-advanced/internal/handler"
	"github.com/38888/nunu-layout-advanced/internal/repository"
	"github.com/38888/nunu-layout-advanced/internal/server"
	"github.com/38888/nunu-layout-advanced/internal/service"
	"github.com/38888/nunu-layout-advanced/pkg/app"
	"github.com/38888/nunu-layout-advanced/pkg/jwt"
	"github.com/38888/nunu-layout-advanced/pkg/log"
	"github.com/38888/nunu-layout-advanced/pkg/server/http"
	"github.com/38888/nunu-layout-advanced/pkg/sid"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	//repository.NewRedis,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewUserRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
)

var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewJob,
)

// build App
func newApp(
	httpServer *http.Server,
	job *server.Job,
	// task *server.Task,
) *app.App {
	return app.NewApp(
		app.WithServer(httpServer, job),
		app.WithName("demo-server"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		serverSet,
		sid.NewSid,
		jwt.NewJwt,
		newApp,
	))
}
