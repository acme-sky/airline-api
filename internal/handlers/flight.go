package handlers

import (
	"net/http"

	"github.com/acme-sky/airline-api/internal/models"
	"github.com/acme-sky/airline-api/pkg/db"
	"github.com/gin-gonic/gin"
)

// Handle GET request for `Flight` model.
// It returns a list of flights.
// GetFlights godoc
//
//	@Summary	Get all flights
//	@Schemes
//	@Description	Get all flights
//	@Tags			Flights
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/v1/flights/ [get]
func FlightHandlerGet(c *gin.Context) {
	db, _ := db.GetDb()
	var flights []models.Flight
	db.Preload("DepartaureAirport").Preload("ArrivalAirport").Find(&flights)

	c.JSON(http.StatusOK, gin.H{
		"count": len(flights),
		"data":  &flights,
	})
}

// Handle POST request for `Flight` model.
// Validate JSON input by the request and crate a new flight. Finally returns
// the new created data (after preloading the foreign key objects).
// PostFlights godoc
//
//	@Summary	Create a new flight
//	@Schemes
//	@Description	Create a new flight
//	@Tags			Flights
//	@Accept			json
//	@Produce		json
//	@Success		201
//	@Router			/v1/flights/filter/ [post]
func FlightHandlerPost(c *gin.Context) {
	db, _ := db.GetDb()
	var input models.FlightInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := models.ValidateFlight(db, input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	flight := models.NewFlight(input)
	db.Create(&flight)
	db.Preload("DepartaureAirport").Preload("ArrivalAirport").Find(&flight)

	c.JSON(http.StatusCreated, flight)
}

// Handle GET request for a selected id.
// Returns the flight or a 404 status
// GetFlightById godoc
//
//	@Summary	Get a flight
//	@Schemes
//	@Description	Get a flight
//	@Tags			Flights
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/v1/flights/{flightId}/ [get]
func FlightHandlerGetId(c *gin.Context) {
	db, _ := db.GetDb()
	var flight models.Flight
	if err := db.Where("id = ?", c.Param("id")).Preload("DepartaureAirport").Preload("ArrivalAirport").First(&flight).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, flight)
}

// Handle PUT request for `Flight` model.
// First checks if the selected flight exists or not. Then, validates JSON input by the
// request and edit a selected flight. Finally returns the new created data.
// EditFlightById godoc
//
//	@Summary	Edit a flight
//	@Schemes
//	@Description	Edit a flight
//	@Tags			Flights
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/v1/flights/{flightId}/ [put]
func FlightHandlerPut(c *gin.Context) {
	db, _ := db.GetDb()
	var flight models.Flight
	if err := db.Where("id = ?", c.Param("id")).First(&flight).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	var input models.FlightInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := models.ValidateFlight(db, input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	db.Model(&flight).Updates(input)
	db.Preload("DepartaureAirport").Preload("ArrivalAirport").Find(&flight)

	c.JSON(http.StatusOK, flight)
}

// Filter flights by departaure (airport and time) and arrival (airport and
// time). This handler can be called by everyone.
// FilterFlights godoc
//
//	@Summary	Filter flights
//	@Schemes
//	@Description	Filter flights
//	@Tags			Flights
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/v1/flights/{flightId}/ [post]
func FlightHandlerFilter(c *gin.Context) {
	db, _ := db.GetDb()

	var input models.FlightFilterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	airports, err := models.ValidateFlightFilter(db, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var flights []models.Flight
	if err := db.Where("departaure_airport_id = ? AND arrival_airport_id = ? AND departaure_time::date = to_date(?, 'YYYY-MM-DD') AND arrival_time::date = to_date(?, 'YYYY-MM-DD')",
		airports[0], airports[1], input.DepartaureTime.Format("2006-01-02"), input.ArrivalTime.Format("2006-01-02")).Preload("DepartaureAirport").Preload("ArrivalAirport").Find(&flights).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(flights),
		"data":  &flights,
	})
}
