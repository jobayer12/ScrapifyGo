package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func SitemapValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		sitemapURL := c.Query("url")

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
