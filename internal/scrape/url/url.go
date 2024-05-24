package url

import (
	"fmt"
	utils "github.com/jobayer12/ScrapifyGo/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

type Response struct {
	urls []string
}

// UrlScrape godoc
//
//	@Summary		Get the sitemap url list
//	@Description	Return sitemap url list. Example of the sitemap url: https://www.shopify.com/sitemap.xml
//	@Tags			sitemap
//	@Router			/api/v1/url [get]
//	@Param			url	query		string	true	"url"
//	@Response		200	{object}	Response
//	@Produce		application/json
func UrlScrape(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
		c.Abort()
		return
	}

	// Create a new collector.
	collector := colly.NewCollector()

	// Slice to hold the scraped URLs.
	var urls []string

	// OnHTML callback to scrape the URLs.
	collector.OnHTML("html", func(e *colly.HTMLElement) {
		links := e.ChildAttrs("a", "href")
		for _, link := range links {
			if utils.URLValidator(link) {
				urls = append(urls, link)
				fmt.Println(link)
			}
		}
	})

	// Visit the url.
	err := collector.Visit(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to visit the url due to: " + err.Error()})
		c.Abort()
		return
	}

	// Wait until scraping is complete.
	collector.Wait()

	uniqueURL := utils.RemoveDuplicates(urls)
	// Return the scraped URLs as JSON.
	c.JSON(http.StatusOK, gin.H{"urls": uniqueURL})
}
