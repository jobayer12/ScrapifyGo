package search

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/jobayer12/ScrapifyGo/utils"
	"net/http"
)

func ScrapeGoogleSearch(c *gin.Context) {

	response := utils.APIResponse[[]string]{
		Error:  "",
		Status: http.StatusOK,
		Data:   []string{},
	}
	url := c.Query("url")

	if url == "" {
		response.Status = http.StatusBadRequest
		response.Error = "url parameter is required"
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Create a new collector.
	domains := "https://www.google.com"
	collector := colly.NewCollector(colly.AllowURLRevisit(), colly.AllowedDomains(domains))

	collector.OnHTML(".MjjYud", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
	})

	collector.Visit(url)
	collector.Wait()

	c.JSON(http.StatusOK, response)
}
