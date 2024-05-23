package validation

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

// QueryParamsURLValidator is a middleware function that checks if the provided URL query parameter is valid.
func QueryParamsURLValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawURL := c.Query("url")
		if rawURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required"})
			c.Abort()
			return
		}

		if !URLValidator(rawURL) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL parameter"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func URLValidator(rawURL string) bool {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil || !parsedURL.IsAbs() {
		return false
	}
	return true
}
