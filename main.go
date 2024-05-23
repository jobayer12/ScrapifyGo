package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jobayer12/scrape-google-search/docs"
	"github.com/jobayer12/scrape-google-search/email"
	"github.com/jobayer12/scrape-google-search/sitemap"
	"github.com/jobayer12/scrape-google-search/url_scrape"
	"github.com/jobayer12/scrape-google-search/validation"
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
// @host localhost:8080
// @BasePath /
func main() {
	server.GET("/sitemap", validation.SitemapValidator(), sitemap.ScrapeSitemap)
	server.GET("/email", validation.QueryParamsURLValidator(), email.ScrapeEmail)
	server.GET("/url", validation.QueryParamsURLValidator(), url_scrape.UrlScrape)
	log.Fatal(server.Run(":8080"))
}
