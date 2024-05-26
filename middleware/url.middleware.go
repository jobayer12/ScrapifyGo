package middleware

import (
	"github.com/gin-gonic/gin"
	utils "github.com/jobayer12/ScrapifyGo/utils"
	"net/http"
)

// URLMiddleware is a middleware function that checks if the provided URL query parameter is valid.
func URLMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		response := utils.APIResponse[[]string]{
			Error:  "",
			Status: http.StatusOK,
			Data:   make([]string, 0),
		}
		rawURL := c.Query("url")
		if rawURL == "" {
			response.Status = http.StatusBadRequest
			response.Error = "URL parameter is required"
			c.JSON(http.StatusBadRequest, response)
			c.Abort()
			return
		}

		if !utils.URLValidator(rawURL) {
			response.Status = http.StatusBadRequest
			response.Error = "Invalid URL parameter"
			c.JSON(http.StatusBadRequest, response)
			c.Abort()
			return
		}

		c.Next()
	}
}
