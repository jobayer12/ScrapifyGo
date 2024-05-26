package url

import (
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	utils "github.com/jobayer12/ScrapifyGo/utils"
	"net/http"
)

// UrlScrape godoc
//
//	@Summary		Get the url list
//	@Tags			url
//	@Router			/api/v1/url [get]
//	@Param			url	query		string	true	"url"
//	@Response		200	{object} utils.APIResponse[[]string]
//	@Produce		application/json
func UrlScrape(c *gin.Context) {
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
		c.Abort()
		return
	}

	// Create a new collector.
	collector := colly.NewCollector()

	// OnHTML callback to scrape the URLs.
	collector.OnHTML("html", func(e *colly.HTMLElement) {
		links := e.ChildAttrs("a", "href")
		for _, link := range links {
			if utils.URLValidator(link) {
				response.Data = append(response.Data, link)
			}
		}
	})

	// Visit the url.
	err := collector.Visit(url)
	if err != nil {
		response.Status = http.StatusBadRequest
		response.Error = "Failed to visit the url due to: " + err.Error()
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	// Wait until scraping is complete.
	collector.Wait()

	response.Data = utils.RemoveDuplicates(response.Data)
	// Return the scraped URLs as JSON.
	c.JSON(http.StatusOK, response)
}
