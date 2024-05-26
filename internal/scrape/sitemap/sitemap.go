package sitemap

import (
	"encoding/json"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/jobayer12/ScrapifyGo/utils"
	_ "github.com/jobayer12/ScrapifyGo/utils"
	"log"
	"net/http"
)

type URL struct {
	Loc        string `xml:"loc" json:"loc"`
	LastMod    string `xml:"lastmod,omitempty"  json:"lastmod,omitempty"`
	Changefreq string `xml:"changefreq,omitempty" json:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty" json:"priority,omitempty"`
}

type URLSet struct {
	URLs []URL `xml:"url"`
}

// ScrapeSitemap godoc
//
//	@Summary		Get the sitemap url list
//	@Description	Return sitemap url list. Example of the sitemap url: https://www.shopify.com/sitemap.xml
//	@Tags			sitemap
//	@Router			/api/v1/sitemap [get]
//	@Param			url	query	string	true	"url"
//	@Response		200	{object}	_.APIResponse[[]URL]
//	@Produce		application/json
func ScrapeSitemap(c *gin.Context) {
	sitemapURL := c.Query("url")
	if sitemapURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
		return
	}

	// Create a new collector.
	collector := colly.NewCollector(colly.AllowedDomains())

	response := utils.APIResponse[[]URL]{
		Error:  "",
		Status: http.StatusOK,
		Data:   []URL{},
	}

	var urls []URL

	// On XML response, parse the XML.
	collector.OnResponse(func(r *colly.Response) {
		var sitemap URLSet
		err := xml.Unmarshal(r.Body, &sitemap)
		if err != nil {
			log.Printf("Failed to unmarshal XML: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse sitemap XML"})
			return
		}
		urls = sitemap.URLs
	})

	// Visit the sitemap url.
	err := collector.Visit(sitemapURL)
	if err != nil {
		response.Error = "Failed to parse sitemap XML" + err.Error()
		response.Status = http.StatusInternalServerError
		c.JSON(http.StatusInternalServerError, response)
		c.Abort()
		return
	}

	// Wait until scraping is complete.
	collector.Wait()

	// Convert the URLs to JSON.
	jsonData, err := json.MarshalIndent(urls, "", "  ")
	if err != nil {
		response.Error = "Failed to marshal JSON due to " + err.Error()
		response.Status = http.StatusInternalServerError
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
		c.Abort()
		return
	}

	err = json.Unmarshal(jsonData, &response.Data)

	if err != nil {
		response.Error = "Failed to marshal JSON due to " + err.Error()
		response.Status = http.StatusInternalServerError
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response)
}
