package rest

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kolosiv/news-aggr/internal/repository"
)

type RestController interface {
	GetNewsByInterval(c *gin.Context)
	GetTodayNews(c *gin.Context)
}

type restController struct {
	nr repository.NewsRepository
}

func CreateRestController(nr repository.NewsRepository) RestController {
	return &restController{nr: nr}
}

func (rc *restController) GetNewsByInterval(c *gin.Context) {
	start := c.Query("start")
	end := c.Query("end")

	// Парсинг дат
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date"})
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date"})
		return
	}

	news, err := rc.nr.GetNewsByInterval(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the news as JSON
	c.JSON(http.StatusOK, news)
}

func (rc *restController) GetTodayNews(c *gin.Context) {
	news, err := rc.nr.GetTodayNews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, news)
}
