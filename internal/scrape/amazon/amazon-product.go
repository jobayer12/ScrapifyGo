package amazon

import (
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"github.com/jobayer12/ScrapifyGo/utils"
	"log/slog"
	"net/http"
	"strings"
)

type Product struct {
	Title       string   `json:"title"`
	Price       string   `json:"price"`
	Rating      string   `json:"rating"`
	TotalRating string   `json:"totalRating"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
}

// ScrapeAmazonSearch godoc
//
//	@Summary		Get the amazon product details
//	@Tags			amazon
//	@Router			/api/v1/amazon [get]
//	@Param			url	query	string	true	"amazon product url"
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

	collector.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", utils.GetRandomUserAgent())
	})

	collector.OnHTML("div#ppd", func(htmlElement *colly.HTMLElement) {
		product.Title = strings.TrimSpace(htmlElement.ChildText("span#productTitle"))
		product.Price = htmlElement.ChildText("span.priceToPay")
		product.Rating = htmlElement.ChildText("#acrPopover span.a-icon-alt")

		product.Description = strings.Join(strings.Fields(strings.TrimSpace(htmlElement.ChildText("#feature-bullets ul"))), " ")
		rating := strings.Split(htmlElement.ChildText("span#acrCustomerReviewText"), " ")
		if len(rating) == 2 {
			product.TotalRating = rating[0]
		}
	})

	collector.OnHTML("ul", func(htmlElement *colly.HTMLElement) {
		var Images []string
		htmlElement.ForEach("li.imageThumbnail", func(i int, element *colly.HTMLElement) {
			var links = element.ChildAttrs("img", "src")
			for _, link := range links {
				if utils.URLValidator(link) {
					Images = append(Images, link)
				}
			}
			product.Images = Images
		})
	})

	//collector.OnHTML("span#productTitle", func(e *colly.HTMLElement) {
	//	product.Title = strings.TrimSpace(e.Text)
	//})
	//
	//collector.OnHTML("#priceblock_ourprice", func(e *colly.HTMLElement) {
	//	product.Price = e.Text
	//})
	//
	//collector.OnHTML("#acrPopover span.a-icon-alt", func(e *colly.HTMLElement) {
	//	product.Rating = e.Text
	//})
	//

	//collector.OnHTML("#feature-bullets ul", func(e *colly.HTMLElement) {
	//	var descriptions []string
	//	slog.Info(e.Text)
	//	e.ForEach("li", func(i int, element *colly.HTMLElement) {
	//		descriptions = append(descriptions, strings.TrimSpace(element.Text))
	//	})
	//	product.Description = descriptions
	//})
	//
	//collector.OnHTML("span#acrCustomerReviewText", func(htmlElement *colly.HTMLElement) {
	//	ratingText := strings.Split(htmlElement.Text, " ")
	//	product.TotalRating = "0"
	//	if len(ratingText) == 2 {
	//		product.TotalRating = ratingText[0]
	//	}
	//})

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

	response.Data = product

	c.JSON(http.StatusOK, response)
}
