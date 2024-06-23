package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jobayer12/ScrapifyGo/docs"
	"github.com/jobayer12/ScrapifyGo/internal/routes"
	"github.com/jobayer12/ScrapifyGo/logger"
	"github.com/jobayer12/ScrapifyGo/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log/slog"
	"strconv"
)

var (
	server       *gin.Engine
	httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_duration_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"path"})

	responseStatus = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_response_status",
		Help: "HTTP response status codes.",
	}, []string{"status_code"})

	totalRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_total_requests",
		Help: "Total number of HTTP requests.",
	}, []string{"path"})
)

func init() {
	prometheus.MustRegister(httpDuration, responseStatus, totalRequests)
	gin.SetMode(gin.ReleaseMode)
	slog.SetDefault(logger.Logger())
	server = gin.Default()
	server.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}

// @title			Scrape API
// @version		1.0
// @description	List of Scrape API
// @host			slight-tiffie-splendid-1fcf1fda.koyeb.app
// @BasePath		/
func main() {
	server.Use(PrometheusMiddleware())
	server.GET("/metrics", gin.WrapH(promhttp.Handler()))
	api := server.Group("api")
	{
		api.Use(middleware.URLMiddleware())
		routes.Routes(api.Group("v1"))
	}

	err := server.Run(":8080")
	if err != nil {
		slog.Error(err.Error())
	}
}

// PrometheusMiddleware is the equivalent middleware
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.FullPath() == "/metrics" {
			c.Next()
			return
		}

		// Get the route path template
		path := c.FullPath()

		// Start the timer
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))

		// Process request
		c.Next()

		// Record status code and total requests
		statusCode := c.Writer.Status()
		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		totalRequests.WithLabelValues(path).Inc()

		// Stop the timer and observe the duration
		timer.ObserveDuration()
	}
}
