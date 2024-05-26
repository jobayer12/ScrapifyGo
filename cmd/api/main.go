package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jobayer12/ScrapifyGo/docs"
	"github.com/jobayer12/ScrapifyGo/internal/routes"
	"github.com/jobayer12/ScrapifyGo/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

var (
	server *gin.Engine
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	server = gin.Default()
	server.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// @title			Kubernetes API
// @version		1.0
// @description	List of kubernetes API
// @host			grotesque-vivianne-splendid-ab71dd99.koyeb.app
// @BasePath		/
func main() {

	api := server.Group("api")
	{
		api.Use(middleware.URLMiddleware())
		routes.Routes(api.Group("v1"))
	}
	log.Fatal(server.Run(":8080"))
}
