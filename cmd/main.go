package main

import (
	"log"

	"github.com/acme-sky/airline-api/internal/handlers"
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
			airports.GET("/", handlers.AirportHandlerGet)
			airports.POST("/", handlers.AirportHandlerPost)
			airports.GET("/:id/", handlers.AirportHandlerGetId)
			airports.PUT("/:id/", handlers.AirportHandlerPut)
		}

		flights := v1.Group("/flights")
		{
			flights.GET("/", handlers.FlightHandlerGet)
			flights.POST("/", handlers.FlightHandlerPost)
			flights.GET("/:id/", handlers.FlightHandlerGetId)
			flights.PUT("/:id/", handlers.FlightHandlerPut)
		}
	}

	router.Run()
}
