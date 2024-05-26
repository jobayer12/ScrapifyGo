package email

import (
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	utils "github.com/jobayer12/ScrapifyGo/utils"
	"net/http"
	"regexp"
)

// ScrapeEmail godoc
//
//	@Summary		Get the email list
//	@Description	Return sitemap url list.
//	@Tags			email
//	@Router			/api/v1/email [get]
//	@Param			url	query	string	true	"url"
//	@Response		200	{object} utils.APIResponse[[]string]
//	@Produce		application/json
func ScrapeEmail(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url parameter is required"})
		return
	}

	// Create a new collector.
	collector := colly.NewCollector(colly.AllowURLRevisit())

	response := utils.APIResponse[[]string]{
		Error:  "",
		Status: http.StatusOK,
		Data:   []string{},
	}

	collector.OnResponse(func(r *colly.Response) {
		response.Data = extractUniqueEmails(string(r.Body))
	})

	// Visit the sitemap url.
	err := collector.Visit(url)
	if err != nil {
		response.Error = "Failed to visit the url due to: " + err.Error()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to visit url"})
		c.Abort()
		return
	}

	// Wait until scraping is complete.
	collector.Wait()

	c.JSON(http.StatusOK, response)
}

func extractUniqueEmails(body string) []string {
	// Define a regular expression pattern to match emails
	emailPattern := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)

	// Find all matches of emails in the body
	emailMatches := emailPattern.FindAllString(body, -1)

	// Remove duplicate emails
	emails := utils.RemoveDuplicates(emailMatches)

	return emails
}
