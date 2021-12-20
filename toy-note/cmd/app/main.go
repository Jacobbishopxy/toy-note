package main

import (
	"toy-note/logger"
	// "github.com/gin-gonic/gin"
	// swaggerFiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
	// _ "gin-swag/docs"
)

const logPath = "../../../logs/toy-note.log"

// @title Toy-note API
// @version 1.0
// @description A simple toy-note API

// @contact.name Jacob Bishop
// @contact.url https://github.com/Jacobbishopxy
// @contact.email jacobbishopxy@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @query.collection.format multi
func main() {
	// Initialize logger
	if err := logger.Init("debug", logPath, false); err != nil {
		panic(err)
	} else {
		defer logger.TNLogger.Sync()
	}

	// Main code here...
	slog := logger.TNLogger.NewSugar("Main")

	slog.Info("Hello world!")

	// router := gin.Default()

	// router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// router.Run()
}
