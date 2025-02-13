package main

import (
	"log"
	"service_api/pkg/handler"
	"service_api/pkg/logger"
	"service_api/pkg/service"

	database "service_api/pkg/db"

	"github.com/joho/godotenv"
)

// @title Service Center API
//@version 1.0
// @description This is a service center API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Author

func main() {
	logger.Log.Info("HTTP Server starting...")

	if err := godotenv.Load(); err != nil {
		logger.Log.Error(".env файл не существует")
	}

	database.DataSource.InitAdmin()

	// Routs
	routes := handler.Handler{
		AuthService: service.AuthService{},
	}

	// HTTP-Server init
	srv := Server{}

	if err := srv.Start("8080", routes.InitRoutes()); err != nil {
		log.Fatal(err)
	}
}
