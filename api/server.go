package api

import (
	"github.com/gin-gonic/gin"
	"github.com/raymondgitonga/love_and_war/engine"
)

const (
	serverAddress = "0.0.0.0:8081"
)

func NewServer() error {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/attack", engine.Attack)

	return router.Run(serverAddress)
}
