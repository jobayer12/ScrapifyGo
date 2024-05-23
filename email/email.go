package email

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"log"
	"net/http"
)

// ScrapeEmail godoc
// @Summary			Get the email list
// @Description		Return sitemap url list.
// @Tags			email
// @Router			/email [get]
// @Param 			url query string true "url"
// @Response		200 {array} string
// @Produce			application/json
func ScrapeEmail(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url_scrape parameter is required"})
		return
	}

	// Create a new collector.
	collector := colly.NewCollector(colly.AllowedDomains())

	var emails []string

	// Error handling.
	collector.OnError(func(_ *colly.Response, err error) {
		log.Printf("Request failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scrape emails"})
		c.Abort()
		return
	})

	collector.OnResponse(func(r *colly.Response) {

		fmt.Println(r.Body)
	})

	collector.OnScraped(func(r *colly.Response) {

	})

	// Visit the sitemap url_scrape.
	err := collector.Visit(url)
	if err != nil {
		log.Printf("Failed to visit url_scrape: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to visit url_scrape"})
		c.Abort()
		return
	}

	// Wait until scraping is complete.
	collector.Wait()

	c.JSON(http.StatusOK, emails)
}
