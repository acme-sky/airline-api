package main

import (
	"log"

	"github.com/acme-sky/airline-api/internal/handlers"
	"github.com/acme-sky/airline-api/pkg/config"
	"github.com/acme-sky/airline-api/pkg/db"
	"github.com/acme-sky/airline-api/pkg/middleware"
	"github.com/gin-gonic/gin"
)

// Create a new instance of Gin server
func main() {
	router := gin.Default()

	// Mode is setted as Debug by default. This line below keep this thing more
	// expressive but it should be setted up by an environment variable.
	// TODO: set mode by env var
	gin.SetMode(gin.DebugMode)

	var err error

	// Read environment variables and stops execution if any errors occur
	err = config.LoadConfig()
	if err != nil {
		log.Printf("failed to load config. err %v", err)

		return
	}

	// Ignore error because if it failed on loading, it should raised an error
	// above.
	config, _ := config.GetConfig()

	if _, err := db.InitDb(config.String("database.dsn")); err != nil {
		log.Printf("failed to connect database. err %v", err)

		return
	}

	// v1 is just like a namespace for every routing here below
	v1 := router.Group("/v1")
	{
		v1.POST("/login/", handlers.LoginHandler)
		airports := v1.Group("/airports")
		{
			airports.Use(middleware.Auth())
			airports.GET("/", handlers.AirportHandlerGet)
			airports.POST("/", handlers.AirportHandlerPost)
			airports.GET("/:id/", handlers.AirportHandlerGetId)
			airports.PUT("/:id/", handlers.AirportHandlerPut)
		}

		flights := v1.Group("/flights")
		{
			flights.GET("/", middleware.Auth(), handlers.FlightHandlerGet)
			flights.POST("/", middleware.Auth(), handlers.FlightHandlerPost)
			flights.GET("/:id/", middleware.Auth(), handlers.FlightHandlerGetId)
			flights.PUT("/:id/", middleware.Auth(), handlers.FlightHandlerPut)
			flights.POST("/filter/", handlers.FlightHandlerFilter)
		}
	}

	router.Run()
}
