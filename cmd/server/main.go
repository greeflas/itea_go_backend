package main

import (
	"log"

	"github.com/greeflas/itea_go_backend/internal/handler"
	"github.com/greeflas/itea_go_backend/pkg/middleware"
	"github.com/greeflas/itea_go_backend/pkg/server"
)

func main() {
	logger := log.Default()

	authMiddleware := middleware.NewAuthMiddleware("secret_token")

	userHandler := handler.NewUserHandler(logger)

	apiServer := server.NewAPIServer(logger)
	apiServer.AddRoute("/user", authMiddleware.Wrap(userHandler))

	if err := apiServer.Start(); err != nil {
		logger.Fatal(err)
	}
}
