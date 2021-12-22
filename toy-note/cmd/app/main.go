package main

import (
	"flag"
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
const envPath = "../../env"

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

	// Determine the api environment
	mode := flag.String("m", "dev", "dev or prod")
	flag.Parse()

	var logLevel string

	// set the log level and gin mode
	if *mode == "dev" {
		logLevel = "debug"
		gin.SetMode(gin.DebugMode)
	} else {
		logLevel = "info"
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize logger
	if err := logger.Init(logLevel, logPath, false); err != nil {
		panic(err)
	} else {
		defer logger.TNLogger.Sync()
	}

	// Main logger
	log := logger.TNLogger.NewSugar("main")
	log.Info("Starting toy-note services...")

	// Load config
	config, err := util.LoadConfig(false, envPath)
	if err != nil {
		log.Panic(err)
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
		log.Panic(err)
	}

	// Initialize controller
	toyNoteController := controller.NewToyNoteController(logger.TNLogger, toyNoteService)
	if err != nil {
		log.Panic(err)
	}

	// Gin
	router := gin.New()

	// Api group
	api := router.Group("/api")
	{
		api.GET("/get-tags", toyNoteController.GetTags)
		api.POST("/save-tag", toyNoteController.SaveTag)
		api.DELETE("/delete-tag/:id", toyNoteController.DeleteTag)

		api.GET("/get-posts", toyNoteController.GetPosts)
		api.POST("/save-post", toyNoteController.SavePost)
		api.DELETE("/delete-post/:id", toyNoteController.DeletePost)

		api.GET("/download-file/:id", toyNoteController.DownloadAffiliate)

		api.GET("/search-posts-by-tags", toyNoteController.SearchPostsByTags)
	}

	// Swagger documention
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	log.Info("Starting toy-note API...")
	router.Run()
}
