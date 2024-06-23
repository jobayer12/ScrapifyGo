package search

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"github.com/jobayer12/ScrapifyGo/utils"
	"log/slog"
	"net/http"
)

type DuckDuckGoSearchResponse struct {
	Url             string   `json:"url"`
	Title           string   `json:"title"`
	MetaDescription string   `json:"metaDescription"`
	RelatedSearch   []string `json:"relatedSearch"`
}

func DuckDuckSearchResult(c *gin.Context) {
	response := utils.APIResponse[[]DuckDuckGoSearchResponse]{
		Error:  "",
		Status: http.StatusOK,
		Data:   []DuckDuckGoSearchResponse{},
	}
	url := c.Query("url")

	if url == "" {
		response.Status = http.StatusBadRequest
		response.Error = "url parameter is required"
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	collector := colly.NewCollector(colly.AllowURLRevisit(), colly.AllowedDomains("duckduckgo.com", "https://duckduckgo.com"))

	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", utils.GetRandomUserAgent())
	})

	collector.OnHTML("ol.react-results--main", func(htmlElement *colly.HTMLElement) {
		htmlElement.ForEach("li[data-layout=organic]", func(i int, element *colly.HTMLElement) {
			fmt.Println(element)
		})
	})
	err := collector.Visit(url)

	collector.Wait()

	if err != nil {
		slog.Error(err.Error())
		response.Status = http.StatusBadRequest
		response.Error = "Failed to collect product details due to " + err.Error()
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	//response.Data = product

	c.JSON(http.StatusOK, response)
}
