package handlers

import (
	"net/http"

	"github.com/acme-sky/airline-api/internal/models"
	"github.com/acme-sky/airline-api/pkg/db"
	"github.com/gin-gonic/gin"
)

// Handle GET request for `Journey` model.
// It returns a list of journeys.
// GetJourney godoc
//
//	@Summary	Get all journeys
//	@Schemes
//	@Description	Get all journeys
//	@Tags			Journeys
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/v1/journeys/ [get]
func JourneyHandlerGet(c *gin.Context) {
	db, _ := db.GetDb()
	var journeys []models.Journey

	db = models.JourneyPreload(db)
	db.Find(&journeys)

	c.JSON(http.StatusOK, gin.H{
		"count": len(journeys),
		"data":  &journeys,
	})
}

// Handle POST request for `Journey` model.
// Validate JSON input by the request and crate a new journey. Finally returns
// the new created data (after preloading the foreign key objects).
// PostJourney godoc
//
//	@Summary	Create a new journey
//	@Schemes
//	@Description	Create a new journey
//	@Tags			Journeys
//	@Accept			json
//	@Produce		json
//	@Success		201
//	@Router			/v1/journeys/ [post]
func JourneyHandlerPost(c *gin.Context) {
	db, _ := db.GetDb()
	var input models.JourneyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := models.ValidateJourney(db, input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	journey := models.NewJourney(input)
	db.Create(&journey)
	db = models.JourneyPreload(db)
	db.Find(&journey)

	c.JSON(http.StatusCreated, journey)
}

// Handle GET request for a selected id.
// Returns the journey or a 404 status
// GetJourneyById godoc
//
//	@Summary	Get a journey
//	@Schemes
//	@Description	Get a journey
//	@Tags			Journeys
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/v1/journeys/{journeyId}/ [get]
func JourneyHandlerGetId(c *gin.Context) {
	db, _ := db.GetDb()
	var journey models.Journey
	db = db.Where("id = ?", c.Param("id"))
	db = models.JourneyPreload(db)
	if err := db.First(&journey).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, journey)
}

// Handle PUT request for `Journey` model.
// First checks if the selected journey exists or not. Then, validates JSON input by the
// request and edit a selected journey. Finally returns the new created data.
// EditJourneyById godoc
//
//	@Summary	Edit a journey
//	@Schemes
//	@Description	Edit a journey
//	@Tags			Journeys
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/v1/journeys/{journeyId}/ [put]
func JourneyHandlerPut(c *gin.Context) {
	db, _ := db.GetDb()
	var journey models.Journey
	if err := db.Where("id = ?", c.Param("id")).First(&journey).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	var input models.JourneyInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := models.ValidateJourney(db, input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	db.Model(&journey).Updates(input)
	db = models.JourneyPreload(db)
	db.Find(&journey)

	c.JSON(http.StatusOK, journey)
}
