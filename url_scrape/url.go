package url_scrape

import (
	"fmt"
	"github.com/jobayer12/scrape-google-search/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

// UrlScrape scrapes all the URLs from the given page and returns them as JSON.
func UrlScrape(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url_scrape parameter is required"})
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
			if validation.URLValidator(link) {
				urls = append(urls, link)
				fmt.Println(link)
			}
		}
	})

	//// Error handling.
	//collector.OnError(func(_ *colly.Response, err error) {
	//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scrape the url_scrape", "details": err.Error()})
	//	c.Abort()
	//	return
	//})

	// Visit the url_scrape.
	err := collector.Visit(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to visit the url_scrape due to: " + err.Error()})
		c.Abort()
		return
	}

	// Wait until scraping is complete.
	collector.Wait()

	// Return the scraped URLs as JSON.
	c.JSON(http.StatusOK, gin.H{"urls": urls})
}
