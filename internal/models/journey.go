package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Journey model
// Save again the cost because it could be changed meanwhile
type Journey struct {
	Id                uint      `gorm:"column:id" json:"id"`
	CreatedAt         time.Time `gorm:"column:created_at" json:"crated_at"`
	DepartureFlightId int       `json:"-"`
	DepartureFlight   Flight    `gorm:"foreignKey:DepartureFlightId" json:"departure_flight"`
	ArrivalFlightId   *int      `gorm:"null" json:"-"`
	ArrivalFlight     *Flight   `gorm:"foreignKey:ArrivalFlightId;null" json:"arrival_flight"`
	Cost              float32   `gorm:"column:cost" json:"cost"`
	Email             string    `gorm:"column:email" json:"email"`
}

func JourneyPreload(db *gorm.DB) *gorm.DB {
	preloads := []string{"DepartureFlight", "ArrivalFlight", "DepartureFlight.DepartureAirport", "DepartureFlight.ArrivalAirport", "ArrivalFlight.DepartureAirport", "ArrivalFlight.ArrivalAirport"}

	for _, preload := range preloads {
		db = db.Preload(preload)
	}

	return db
}

// Struct used to get new data for a journey
type JourneyInput struct {
	DepartureFlightId int     `json:"departure_flight_id" binding:"required"`
	ArrivalFlightId   *int    `json:"arrival_flight_id"`
	Cost              float32 `json:"cost" binding:"required"`
	Email             string  `json:"email" binding:"required"`
}

// It validates data from `in` and returns a possible error or not
func ValidateJourney(db *gorm.DB, in JourneyInput) error {
	var departure_flight Flight
	departure_flight_id := in.DepartureFlightId
	arrival_flight_id := in.ArrivalFlightId
	if err := db.Where("id = ?", departure_flight_id).First(&departure_flight).Error; err != nil {
		return errors.New("`departure_flight_id` does not exist.")
	}

	if arrival_flight_id != nil {
		var arrival_flight Flight
		if err := db.Where("id = ?", arrival_flight_id).First(&arrival_flight).Error; err != nil {
			return errors.New("`arrival_flight_id` does not exist.")
		}

		if departure_flight_id == *arrival_flight_id {
			return errors.New("`departure_flight_id` can't be equals to `arrival_flight_id`")
		}
	}

	return nil
}

// Returns a new Journey with the data from `in`. It should be called after
// `ValidateJourney(..., in)` method
func NewJourney(in JourneyInput) Journey {
	return Journey{
		CreatedAt:         time.Now(),
		DepartureFlightId: in.DepartureFlightId,
		ArrivalFlightId:   in.ArrivalFlightId,
		Cost:              in.Cost,
		Email:             in.Email,
	}
}
