package main

import (
	"toy-note/api/controller"
	"toy-note/api/persistence"
	"toy-note/api/service"
	"toy-note/api/util"
	_ "toy-note/docs"
	"toy-note/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const logPath = "../../../logs/toy-note.log"
const envPath = "../../../env"

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

	// Load config
	config, err := util.LoadConfig(false, envPath)
	if err != nil {
		panic(err)
	}

	// Pg
	pgConn := persistence.PgConn{
		Host:    config.PG_HOST,
		Port:    config.PG_PORT,
		User:    config.PG_USER,
		Pass:    config.PG_PASS,
		Db:      config.PG_DB,
		Sslmode: "disable",
	}

	// Mongo
	mongoConn := persistence.MongoConn{
		Host: config.MONGO_HOST,
		Port: config.MONGO_PORT,
		User: config.MONGO_USER,
		Pass: config.MONGO_PASS,
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

	// Gin
	router := gin.New()

	// Api group
	api := router.Group("/api")
	{
		api.GET("/get-tags", toyNoteController.GetTags)
		api.POST("/save-tag", toyNoteController.SaveTag)
		api.DELETE("/delete-tag", toyNoteController.DeleteTag)
	}

	// Swagger documention
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	router.Run()
}
