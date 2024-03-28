package main

import (
	"log"

	"github.com/acme-sky/airline-api/internal/handlers"
	"github.com/acme-sky/airline-api/pkg/config"
	"github.com/acme-sky/airline-api/pkg/db"
	"github.com/gin-gonic/gin"
)

// Create a new instance of Gin server
func main() {
	router := gin.Default()

	// Mode is setted as Debug by default. This line below keep this thing more
	// expressive but it should be setted up by an environment variable.
	// TODO: set mode by env var
	gin.SetMode(gin.DebugMode)

	// Read environment variables and stops execution if any errors occur
	if err := config.LoadConfig(); err != nil {
		log.Printf("failed to load config. err %v", err)

		return
	}

	_, err := db.InitDb(config.GetConfig().String("database.dsn"))
	if err != nil {
		log.Printf("failed to connect database. err %v", err)

		return
	}

	// v1 is just like a namespace for every routing here below
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
