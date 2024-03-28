package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Flight model
type Flight struct {
	Id                  uint      `gorm:"column:id" json:"id"`
	CreatedAt           time.Time `gorm:"column:created_at" json:"crated_at"`
	DepartaureTime      time.Time `gorm:"column:departaure_time" json:"departaure_time"`
	ArrivalTime         time.Time `gorm:"column:arrival_time" json:"arrival_time"`
	DepartaureAirportId int       `json:"-"`
	DepartaureAirport   Airport   `gorm:"foreignKey:DepartaureAirportId" json:"departaure_airport"`
	ArrivalAirportId    int       `json:"-"`
	ArrivalAirport      Airport   `gorm:"foreignKey:ArrivalAirportId" json:"arrival_airport"`
	Cost                float32   `gorm:"column:cost" json:"cost"`
}

// Struct used to get new data for a flight
type FlightInput struct {
	DepartaureTime      time.Time `json:"departaure_time" binding:"required"`
	ArrivalTime         time.Time `json:"arrival_time" binding:"required"`
	DepartaureAirportId int       `json:"departaure_airport_id" binding:"required"`
	ArrivalAirportId    int       `json:"arrival_airport_id" binding:"required"`
	Cost                float32   `json:"cost" binding:"required"`
}

// It validates data from `in` and returns a possible error or not
func ValidateFlight(db *gorm.DB, in FlightInput) error {
	var departaure_airport Airport
	if err := db.Where("id = ?", in.DepartaureAirportId).First(&departaure_airport).Error; err != nil {
		return errors.New("`departaure_airport_id` does not exist.")
	}

	var arrival_airport Airport
	if err := db.Where("id = ?", in.ArrivalAirportId).First(&arrival_airport).Error; err != nil {
		return errors.New("`arrival_airport_id` does not exist.")
	}

	if in.DepartaureAirportId == in.ArrivalAirportId {
		return errors.New("`departaure_airport_id` can't be equals to `arrival_airport_id`")
	}

	if in.DepartaureTime.Equal(in.ArrivalTime) || in.DepartaureTime.After(in.ArrivalTime) {
		return errors.New("`departaure_time` can't be after or the same `arrival_time`")
	}

	return nil
}

// Returns a new Flight with the data from `in`. It should be called after
// `ValidateFlight(..., in)` method
func NewFlight(in FlightInput) Flight {
	return Flight{
		CreatedAt:           time.Now(),
		DepartaureTime:      in.DepartaureTime,
		DepartaureAirportId: in.DepartaureAirportId,
		ArrivalTime:         in.ArrivalTime,
		ArrivalAirportId:    in.ArrivalAirportId,
		Cost:                in.Cost,
	}
}
