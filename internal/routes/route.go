package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jobayer12/ScrapifyGo/internal/scrape/email"
	"github.com/jobayer12/ScrapifyGo/internal/scrape/sitemap"
	"github.com/jobayer12/ScrapifyGo/internal/scrape/url"
	"github.com/jobayer12/ScrapifyGo/middleware"
)

func Routes(g *gin.RouterGroup) {
	g.GET("/sitemap", middleware.SitemapValidator(), sitemap.ScrapeSitemap)
	g.GET("/email", email.ScrapeEmail)
	g.GET("/url", url.UrlScrape)
}
