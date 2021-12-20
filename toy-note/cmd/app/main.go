package main

import (
	"toy-note/api/controller"
	"toy-note/api/persistence"
	"toy-note/api/service"
	_ "toy-note/docs"
	"toy-note/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const logPath = "../../../logs/toy-note.log"

// @title                    Toy-note API
// @version                  1.0
// @description              A simple toy-note API
// @contact.name             Jacob Bishop
// @contact.url              https://github.com/Jacobbishopxy
// @contact.email            jacobbishopxy@gmail.com
// @license.name             Apache 2.0
// @license.url              http://www.apache.org/licenses/LICENSE-2.0
// @host                     localhost:8080
// @BasePath                 /api
// @query.collection.format  multi
func main() {
	// Initialize logger
	if err := logger.Init("debug", logPath, false); err != nil {
		panic(err)
	} else {
		defer logger.TNLogger.Sync()
	}

	pgConn := persistence.PgConn{
		// TODO:
	}

	mongoConn := persistence.MongoConn{
		// TODO:
	}

	// Initialize service
	toyNoteService, err := service.NewToyNoteService(logger.TNLogger, pgConn, mongoConn)
	if err != nil {
		panic(err)
	}

	// Initialize controller
	toyNoteController := controller.NewToyNoteController(logger.TNLogger, toyNoteService)
	if err != nil {
		panic(err)
	}

	router := gin.New()
	api := router.Group("/api")

	api.GET("/get-tags", toyNoteController.GetTags)
	api.POST("/save-tag", toyNoteController.SaveTag)
	api.DELETE("/delete-tag", toyNoteController.DeleteTag)

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run()
}
