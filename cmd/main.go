package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	gin.SetMode(gin.DebugMode)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"hello": "world",
		})
	})

	router.Run()
}
