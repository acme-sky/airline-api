package handlers

import (
	"net/http"

	"github.com/acme-sky/airline-api/internal/models"
	"github.com/acme-sky/airline-api/pkg/db"
	"github.com/gin-gonic/gin"
)

// Handle GET request for `Airport` model.
// It returns a list of airports.
func AirportHandlerGet(c *gin.Context) {
	db, _ := db.GetDb()

	var airports []models.Airport
	db.Find(&airports)

	c.JSON(http.StatusOK, gin.H{
		"count": len(airports),
		"data":  &airports,
	})
}

// Handle POST request for `Airport` model.
// Validate JSON input by the request and crate a new airport. Finally returns
// the new created data.
func AirportHandlerPost(c *gin.Context) {
	db, _ := db.GetDb()
	var input models.AirportInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	airport := models.NewAirport(input)
	db.Create(&airport)

	c.JSON(http.StatusCreated, airport)
}

// Handle GET request for a selected id.
// Returns the aiport or a 404 status
func AirportHandlerGetId(c *gin.Context) {
	db, _ := db.GetDb()
	var airport models.Airport
	if err := db.Where("id = ?", c.Param("id")).First(&airport).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, airport)
}

// Handle PUT request for `Airport` model.
// First checks if the selected airport exists or not. Then, validates JSON input by the
// request and edit a selected airport. Finally returns the new created data.
func AirportHandlerPut(c *gin.Context) {
	db, _ := db.GetDb()
	var airport models.Airport
	if err := db.Where("id = ?", c.Param("id")).First(&airport).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	var input models.AirportInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&airport).Updates(input)

	c.JSON(http.StatusOK, airport)
}
