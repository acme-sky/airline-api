package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/acme-sky/airline-api/internal/models"
	"github.com/acme-sky/airline-api/pkg/db"
	"github.com/gin-gonic/gin"
)

// Handle GET request for `Hook` model.
// It returns a list of hooks.
func HookHandlerGet(c *gin.Context) {
	db, _ := db.GetDb()

	var hooks []models.Hook
	db.Find(&hooks)

	c.JSON(http.StatusOK, gin.H{
		"count": len(hooks),
		"data":  &hooks,
	})
}

// Handle POST request for `Hook` model.
// Validate JSON input by the request and crate a new hook. Finally returns
// the new created data.
func HookHandlerPost(c *gin.Context) {
	db, _ := db.GetDb()
	var input models.HookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	hook := models.NewHook(input)
	db.Create(&hook)

	c.JSON(http.StatusCreated, hook)
}

// Handle GET request for a selected id.
// Returns the hook or a 404 status
func HookHandlerGetId(c *gin.Context) {
	db, _ := db.GetDb()
	var hook models.Hook
	if err := db.Where("id = ?", c.Param("id")).First(&hook).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hook)
}

// Handle PUT request for `Hook` model.
// First checks if the selected hook exists or not. Then, validates JSON input by the
// request and edit a selected hook. Finally returns the new created data.
func HookHandlerPut(c *gin.Context) {
	db, _ := db.GetDb()
	var hook models.Hook
	if err := db.Where("id = ?", c.Param("id")).First(&hook).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	var input models.HookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&hook).Updates(input)

	c.JSON(http.StatusOK, hook)
}

// Handle POST request to send an offer to all the saved hooks.
// First get all hooks, then validate the request payload which must be
// `{"flight_id": <valid_id>}`
// and finally send the flight object to all the hooks by their endpoint.
func HookHandlerOffer(c *gin.Context) {
	db, _ := db.GetDb()

	var hooks []models.Hook
	db.Find(&hooks)

	var input models.OfferInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var flight models.Flight
	if err := db.Where("id = ?", input.FlightId).Preload("DepartaureAirport").Preload("ArrivalAirport").First(&flight).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	total := 0

	payload, err := json.Marshal(flight)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	reader := bytes.NewReader(payload)

	for _, hook := range hooks {
		req, err := http.NewRequest(http.MethodPost, hook.Endpoint, reader)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			log.Printf("received an error for hook `%d`: %s\n", hook.Id, err.Error())

			continue
		}

		client := http.Client{Timeout: time.Minute}
		res, err := client.Do(req)

		if err != nil || res == nil {
			if err != nil {
				log.Printf("received an error for hook `%d`: %s\n", hook.Id, err.Error())
			}

			if res != nil {
				log.Printf("received an error for hook `%d` with status code %d\n", hook.Id, res.StatusCode)
			}

			continue
		}

		total += 1
	}

	// Send back some info just to know how many hook work or not
	c.JSON(http.StatusOK, gin.H{
		"hooks":  len(hooks),
		"sent":   total,
		"errors": len(hooks) - total,
	})
}
