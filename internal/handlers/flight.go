package handlers

import (
	"net/http"

	"github.com/acme-sky/airline-api/internal/models"
	"github.com/acme-sky/airline-api/pkg/db"
	"github.com/gin-gonic/gin"
)

func FlightHandlerGet(ctx *gin.Context) {
	db := db.GetDb()
	var flights []models.Flight
	db.Preload("DepartaureAirport").Preload("ArrivalAirport").Find(&flights)

	ctx.JSON(http.StatusOK, gin.H{
		"count": len(flights),
		"data":  &flights,
	})
}

func FlightHandlerPost(ctx *gin.Context) {
	db := db.GetDb()
	var input models.FlightInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var departaure_airport models.Airport
	if err := db.Where("id = ?", input.DepartaureAirportId).First(&departaure_airport).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "`departaure_airport_id` does not exist."})
		return
	}
	var arrival_airport models.Airport
	if err := db.Where("id = ?", input.ArrivalAirportId).First(&arrival_airport).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "`arrival_airport_id` does not exist."})
		return
	}

	flight, err := models.NewFlight(input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&flight)
	db.Preload("DepartaureAirport").Preload("ArrivalAirport").Find(&flight)

	ctx.JSON(http.StatusCreated, flight)
}

func FlightHandlerGetId(ctx *gin.Context) {
	db := db.GetDb()
	var flight models.Flight
	if err := db.Where("id = ?", ctx.Param("id")).First(&flight).Error; err != nil {
		ctx.JSON(http.StatusNotFound, map[string]string{})
		return
	}

	ctx.JSON(http.StatusOK, flight)
}

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

	var departaure_airport models.Airport
	if err := db.Where("id = ?", input.DepartaureAirportId).First(&departaure_airport).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "`departaure_airport_id` does not exist."})
		return
	}
	var arrival_airport models.Airport
	if err := db.Where("id = ?", input.ArrivalAirportId).First(&arrival_airport).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "`arrival_airport_id` does not exist."})
		return
	}

	if err := models.ValidateFlight(input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&flight).Updates(input)
	db.Preload("DepartaureAirport").Preload("ArrivalAirport").Find(&flight)

	ctx.JSON(http.StatusOK, flight)
}
