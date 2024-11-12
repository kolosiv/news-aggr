package rest

import (
	"github.com/gin-gonic/gin"
)

func Setup(rc RestController) {
	r := gin.Default()

	// r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP)
	// r.Use(middleware.Recoverer)
	// r.Use(middleware.Timeout(timeout))

	r.GET("/news", rc.GetTodayNews)
	r.GET("/news/interval", rc.GetNewsByInterval)
	r.Run()
}
