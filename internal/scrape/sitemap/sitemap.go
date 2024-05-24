package sitemap

import (
	"encoding/json"
	"encoding/xml"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"log"
	"net/http"
	"strings"
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

func ValidateSitemapURL() gin.HandlerFunc {
	return func(c *gin.Context) {
		sitemapURL := c.Query("url")
		if sitemapURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
			c.Abort()
			return
		}

		// Check if url ends with "sitemap.xml".
		if !strings.HasSuffix(sitemapURL, "sitemap.xml") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "url must end with sitemap.xml"})
			c.Abort()
			return
		}

		// Send a HEAD request to check the content type.
		resp, err := http.Head(sitemapURL)
		if err != nil {
			log.Printf("HEAD request failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify sitemap url"})
			c.Abort()
			return
		}
		defer resp.Body.Close()

		// Check if the content type is XML.
		contentType := resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "xml") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "url does not point to a valid sitemap (invalid content type)"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// ScrapeSitemap godoc
//	@Summary		Get the sitemap url list
//	@Description	Return sitemap url list. Example of the sitemap url: https://www.shopify.com/sitemap.xml
//	@Tags			sitemap
//	@Router			/api/v1/sitemap [get]
//	@Param			url	query	string	true	"url"
//	@Response		200	{array}	URL
//	@Produce		application/json
func ScrapeSitemap(c *gin.Context) {
	sitemapURL := c.Query("url")
	if sitemapURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
		return
	}

	// Create a new collector.
	collector := colly.NewCollector(colly.AllowedDomains())

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
		log.Printf("Failed to visit url: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to visit sitemap url"})
		c.Abort()
		return
	}

	// Wait until scraping is complete.
	collector.Wait()

	// Convert the URLs to JSON.
	jsonData, err := json.MarshalIndent(urls, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal JSON: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON"})
		c.Abort()
		return
	}

	// Return the JSON data.
	c.Data(http.StatusOK, "application/json", jsonData)
}
