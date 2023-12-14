package main

import (
	"github.com/greeflas/itea_go_backend/internal/repository"
	"github.com/greeflas/itea_go_backend/internal/service"
	"go.uber.org/fx"
	"log"

	"github.com/greeflas/itea_go_backend/internal/handler"
	"github.com/greeflas/itea_go_backend/pkg/middleware"
	"github.com/greeflas/itea_go_backend/pkg/server"
)

func main() {
	options := []fx.Option{
		fx.Provide(
			repository.NewUserInMemoryRepository,
			service.NewUserService,
			handler.NewUserHandler,
			server.NewAPIServer,
		),
		fx.Provide(func() *middleware.AuthMiddleware {
			return middleware.NewAuthMiddleware("secret_token")
		}),
		fx.Provide(func() *log.Logger {
			return log.Default()
		}),
		fx.Invoke(func(
			apiServer *server.APIServer,
			authMiddleware *middleware.AuthMiddleware,
			userHandler *handler.UserHandler,
		) {
			apiServer.AddRoute("/user", authMiddleware.Wrap(userHandler))
		}),
		fx.Invoke(func(apiServer *server.APIServer, logger *log.Logger) {
			if err := apiServer.Start(); err != nil {
				logger.Fatal(err)
			}
		}),
	}

	fx.New(options...).Run()
}
