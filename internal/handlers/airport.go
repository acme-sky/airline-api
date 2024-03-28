package handlers

import (
	"net/http"

	"github.com/acme-sky/airline-api/internal/models"
	"github.com/acme-sky/airline-api/pkg/db"
	"github.com/gin-gonic/gin"
)

// Handle GET request for `Airport` model.
// It returns a list of airports.
func AirportHandlerGet(ctx *gin.Context) {
	db := db.GetDb()

	var airports []models.Airport
	db.Find(&airports)

	ctx.JSON(http.StatusOK, gin.H{
		"count": len(airports),
		"data":  &airports,
	})
}

// Handle POST request for `Airport` model.
// Validate JSON input by the request and crate a new airport. Finally returns
// the new created data.
func AirportHandlerPost(ctx *gin.Context) {
	db := db.GetDb()
	var input models.AirportInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	airport := models.NewAirport(input)
	db.Create(&airport)

	ctx.JSON(http.StatusCreated, airport)
}

// Handle GET request for a selected id.
// Returns the aiport or a 404 status
func AirportHandlerGetId(ctx *gin.Context) {
	db := db.GetDb()
	var airport models.Airport
	if err := db.Where("id = ?", ctx.Param("id")).First(&airport).Error; err != nil {
		ctx.JSON(http.StatusNotFound, map[string]string{})
		return
	}

	ctx.JSON(http.StatusOK, airport)
}

// Handle PUT request for `Airport` model.
// First checks if the selected airport exists or not. Then, validates JSON input by the
// request and edit a selected airport. Finally returns the new created data.
func AirportHandlerPut(ctx *gin.Context) {
	db := db.GetDb()
	var airport models.Airport
	if err := db.Where("id = ?", ctx.Param("id")).First(&airport).Error; err != nil {
		ctx.JSON(http.StatusNotFound, map[string]string{})
		return
	}

	var input models.AirportInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&airport).Updates(input)

	ctx.JSON(http.StatusOK, airport)
}
