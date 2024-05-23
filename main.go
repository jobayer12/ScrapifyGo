package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jobayer12/scrape-google-search/docs"
	"github.com/jobayer12/scrape-google-search/sitemap"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

var (
	server *gin.Engine
)

func init() {
	server = gin.Default()
	server.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// @title Kubernetes API
// @version 1.0
// @description List of kubernetes API
// @host grotesque-vivianne-splendid-ab71dd99.koyeb.app
// @BasePath /
func main() {
	server.GET("/sitemap", sitemap.ValidateSitemapURL(), sitemap.ScrapeSitemap)
	log.Fatal(server.Run(":8080"))
}
