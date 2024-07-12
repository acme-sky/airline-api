package main

import (
	"log"

	"github.com/acme-sky/airline-api/internal/handlers"
	"github.com/acme-sky/airline-api/pkg/config"
	"github.com/acme-sky/airline-api/pkg/db"
	"github.com/acme-sky/airline-api/pkg/middleware"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

// Create a new instance of Gin server
func main() {
	router := gin.Default()

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

	// Env variable `debug` set up the mode below
	if !config.Bool("debug") {
		gin.SetMode(gin.ReleaseMode)
	}

	router.Use(cors.AllowAll())

	router.StaticFile("/swagger.yml", "cmd/swagger.yml")

	// v1 is just like a namespace for every routing here below
	v1 := router.Group("/v1")
	{
		v1.POST("/login/", handlers.LoginHandler)
		airports := v1.Group("/airports")
		{
			airports.GET("/", middleware.Auth(), handlers.AirportHandlerGet)
			airports.POST("/", middleware.Auth(), handlers.AirportHandlerPost)
			airports.GET("/:id/", middleware.Auth(), handlers.AirportHandlerGetId)
			airports.GET("/code/:code/", handlers.AirportHandlerGetCode)
			airports.PUT("/:id/", middleware.Auth(), handlers.AirportHandlerPut)
		}

		flights := v1.Group("/flights")
		{
			flights.GET("/", middleware.Auth(), handlers.FlightHandlerGet)
			flights.POST("/", middleware.Auth(), handlers.FlightHandlerPost)
			flights.GET("/:id/", middleware.Auth(), handlers.FlightHandlerGetId)
			flights.PUT("/:id/", middleware.Auth(), handlers.FlightHandlerPut)
			flights.POST("/filter/", handlers.FlightHandlerFilter)
		}

		hooks := v1.Group("/hooks")
		{
			hooks.Use(middleware.Auth())
			hooks.GET("/", handlers.HookHandlerGet)
			hooks.POST("/", handlers.HookHandlerPost)
			hooks.GET("/:id/", handlers.HookHandlerGetId)
			hooks.PUT("/:id/", handlers.HookHandlerPut)
			hooks.POST("/offer/", handlers.HookHandlerOffer)
		}

		journeys := v1.Group("/journeys")
		{
			journeys.Use(middleware.Auth())
			journeys.GET("/", handlers.JourneyHandlerGet)
			journeys.POST("/", handlers.JourneyHandlerPost)
			journeys.GET("/:id/", handlers.JourneyHandlerGetId)
			journeys.PUT("/:id/", handlers.JourneyHandlerPut)
		}
	}

	router.Run()
}
