package google

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/jobayer12/ScrapifyGo/utils"
	"net/http"
)

// ScrapeGoogleSearch godoc
//
//	@Summary		Get the google search list
//	@Tags			google
//	@Router			/api/v1/google [get]
//	@Param			url	query	string	true	"url"
//	@Response		200	{object} utils.APIResponse[[]string]
//	@Produce		application/json
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
	collector := colly.NewCollector(colly.AllowURLRevisit(), colly.AllowedDomains("google.com", "https://www.google.com", "www.google.com"))

	collector.OnHTML(".MjjYud", func(e *colly.HTMLElement) {
		fmt.Println(e.Text)
	})

	//collector.OnResponse(func(r *colly.Response) {
	//	fmt.Println()
	//})

	collector.Visit(url)
	collector.Wait()

	c.JSON(http.StatusOK, response)
}
