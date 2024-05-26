package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/jobayer12/ScrapifyGo/utils"
	"net/http"
	"strings"
)

func SitemapValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := utils.APIResponse[[]string]{
			Error:  "",
			Status: http.StatusOK,
			Data:   make([]string, 0),
		}

		sitemapURL := c.Query("url")

		// Check if url ends with "sitemap.xml".
		if !strings.HasSuffix(sitemapURL, "sitemap.xml") {
			response.Status = http.StatusBadRequest
			response.Error = "url must end with sitemap.xml"
			c.JSON(http.StatusBadRequest, response)
			c.Abort()
			return
		}

		// Send a HEAD request to check the content type.
		resp, err := http.Head(sitemapURL)
		if err != nil {
			response.Status = http.StatusBadRequest
			response.Error = "Failed to verify sitemap url"
			c.JSON(http.StatusInternalServerError, response)
			c.Abort()
			return
		}
		defer resp.Body.Close()

		// Check if the content type is XML.
		contentType := resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "xml") {
			response.Status = http.StatusBadRequest
			response.Error = "url does not point to a valid sitemap (invalid content type)"
			c.JSON(http.StatusBadRequest, response)
			c.Abort()
			return
		}

		c.Next()
	}
}
