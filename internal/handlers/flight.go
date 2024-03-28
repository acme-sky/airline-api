package handlers

import (
	"net/http"

	"github.com/acme-sky/airline-api/internal/models"
	"github.com/acme-sky/airline-api/pkg/db"
	"github.com/gin-gonic/gin"
)

// Handle GET request for `Flight` model.
// It returns a list of flights.
func FlightHandlerGet(ctx *gin.Context) {
	db := db.GetDb()
	var flights []models.Flight
	db.Preload("DepartaureAirport").Preload("ArrivalAirport").Find(&flights)

	ctx.JSON(http.StatusOK, gin.H{
		"count": len(flights),
		"data":  &flights,
	})
}

// Handle POST request for `Flight` model.
// Validate JSON input by the request and crate a new flight. Finally returns
// the new created data (after preloading the foreign key objects).
func FlightHandlerPost(ctx *gin.Context) {
	db := db.GetDb()
	var input models.FlightInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.ValidateFlight(db, input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	flight := models.NewFlight(input)
	db.Create(&flight)
	db.Preload("DepartaureAirport").Preload("ArrivalAirport").Find(&flight)

	ctx.JSON(http.StatusCreated, flight)
}

// Handle GET request for a selected id.
// Returns the flight or a 404 status
func FlightHandlerGetId(ctx *gin.Context) {
	db := db.GetDb()
	var flight models.Flight
	if err := db.Where("id = ?", ctx.Param("id")).Preload("DepartaureAirport").Preload("ArrivalAirport").First(&flight).Error; err != nil {
		ctx.JSON(http.StatusNotFound, map[string]string{})
		return
	}

	ctx.JSON(http.StatusOK, flight)
}

// Handle PUT request for `Flight` model.
// First checks if the selected flight exists or not. Then, validates JSON input by the
// request and edit a selected flight. Finally returns the new created data.
func FlightHandlerPut(ctx *gin.Context) {
	db := db.GetDb()
	var flight models.Flight
	if err := db.Where("id = ?", ctx.Param("id")).First(&flight).Error; err != nil {
		ctx.JSON(http.StatusNotFound, map[string]string{})
		return
	}

	var input models.FlightInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.ValidateFlight(db, input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&flight).Updates(input)
	db.Preload("DepartaureAirport").Preload("ArrivalAirport").Find(&flight)

	ctx.JSON(http.StatusOK, flight)
}
