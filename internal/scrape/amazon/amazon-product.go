package amazon

import (
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"github.com/jobayer12/ScrapifyGo/utils"
	"log"
	"net/http"
	"strings"
)

type Comment struct {
}

type Product struct {
	Title       string   `json:"title"`
	Price       string   `json:"price"`
	Rating      string   `json:"rating"`
	TotalRating string   `json:"totalRating"`
	Description []string `json:"description"`
	Images      []string `json:"images"`
}

// ScrapeAmazonSearch godoc
//
//	@Summary		Get the amazon product details
//	@Tags			amazon
//	@Router			/api/v1/amazon [get]
//	@Param			url	query	string	true	"url"
//	@Response		200	{object} utils.APIResponse[Product]
//	@Produce		application/json
func ScrapeAmazonSearch(c *gin.Context) {
	response := utils.APIResponse[Product]{
		Error:  "",
		Status: http.StatusOK,
		Data:   Product{},
	}
	url := c.Query("url")

	if url == "" {
		response.Status = http.StatusBadRequest
		response.Error = "url parameter is required"
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}
	var product Product

	collector := colly.NewCollector(colly.AllowURLRevisit())

	collector.OnHTML("span#productTitle", func(e *colly.HTMLElement) {
		product.Title = strings.TrimSpace(e.Text)
	})

	collector.OnHTML("#priceblock_ourprice", func(e *colly.HTMLElement) {
		product.Price = e.Text
	})

	collector.OnHTML("#acrPopover span.a-icon-alt", func(e *colly.HTMLElement) {
		product.Rating = e.Text
	})

	collector.OnHTML("#imgTagWrapperId img", func(e *colly.HTMLElement) {
		product.Images = e.ChildAttrs("img", "src")
	})

	collector.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	collector.OnHTML("#feature-bullets ul", func(e *colly.HTMLElement) {
		var descriptions []string
		e.ForEach("li", func(i int, element *colly.HTMLElement) {
			descriptions = append(descriptions, strings.TrimSpace(element.Text))
		})
		product.Description = descriptions
	})

	collector.OnHTML("span#acrCustomerReviewText", func(htmlElement *colly.HTMLElement) {
		ratingText := strings.Split(htmlElement.Text, " ")
		product.TotalRating = "0"
		if len(ratingText) == 2 {
			product.TotalRating = ratingText[0]
		}

	})

	collector.OnHTML("ul.regularAltImageViewLayout", func(htmlElement *colly.HTMLElement) {
		var Images []string
		htmlElement.ForEach("li", func(i int, element *colly.HTMLElement) {
			var links = element.ChildAttrs("img", "src")
			for _, link := range links {
				if utils.URLValidator(link) {
					Images = append(Images, link)
				}
			}
			product.Images = Images
		})
	})

	err := collector.Visit(url)

	collector.Wait()

	if err != nil {
		response.Status = http.StatusBadRequest
		response.Error = "Failed to collect product details"
		c.JSON(http.StatusBadRequest, response)
		c.Abort()
		return
	}

	response.Data = product

	c.JSON(http.StatusOK, response)
}
