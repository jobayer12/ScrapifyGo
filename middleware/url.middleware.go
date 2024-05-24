package middleware

import (
	"github.com/gin-gonic/gin"
	utils "github.com/jobayer12/ScrapifyGo/utils"
	"net/http"
)

// URLMiddleware is a middleware function that checks if the provided URL query parameter is valid.
func URLMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawURL := c.Query("url")
		if rawURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required"})
			c.Abort()
			return
		}

		if !utils.URLValidator(rawURL) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL parameter"})
			c.Abort()
			return
		}

		c.Next()
	}
}
