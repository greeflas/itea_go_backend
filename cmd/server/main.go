package main

import (
	"database/sql"
	"github.com/greeflas/itea_go_backend/internal/repository"
	"github.com/greeflas/itea_go_backend/internal/service"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
	"log"

	"github.com/greeflas/itea_go_backend/internal/handler"
	"github.com/greeflas/itea_go_backend/pkg/middleware"
	"github.com/greeflas/itea_go_backend/pkg/server"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	options := []fx.Option{
		fx.Provide(
			repository.NewUserInMemoryRepository,
			repository.NewUserBunRepository,
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
		fx.Provide(func() *bun.DB {
			dsn := "postgres://postgres:pass@localhost:5432/itea_backend?sslmode=disable"
			sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

			db := bun.NewDB(sqldb, pgdialect.New())

			return db
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
