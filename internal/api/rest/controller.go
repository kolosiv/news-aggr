package rest

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/kolosiv/news-aggr/internal/repository"
)

type RestController interface {
	GetNewsByDate(c *gin.Context)
}

type restController struct {
	nr repository.NewsRepository
}

func CreateRestController(nr repository.NewsRepository) RestController {
	return &restController{nr: nr}
}

func (rc *restController) GetNewsByDate(c *gin.Context) {
	fmt.Printf("HTTP Запрос")
}
