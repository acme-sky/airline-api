package main

import (
	"log"

	"github.com/acme-sky/airline-api/internal/controllers"
	"github.com/acme-sky/airline-api/pkg/config"
	"github.com/acme-sky/airline-api/pkg/db"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	gin.SetMode(gin.DebugMode)

	conf, err := config.LoadConfig()
	if err != nil {
		log.Printf("failed to load config. err %v", err)

		return
	}

	_, err = db.InitDb(conf.String("database.dsn"))
	if err != nil {
		log.Printf("failed to connect database. err %v", err)

		return
	}

	v1 := router.Group("/v1")
	{
		airports := v1.Group("/airports")
		{
			airports.GET("/", controllers.AirportGet)
			airports.POST("/", controllers.AirportPost)
		}
	}

	router.Run()
}
