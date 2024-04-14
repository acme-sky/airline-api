package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Flight model
type Flight struct {
	Id                  uint      `gorm:"column:id" json:"id"`
	Code                string    `gorm:"code" json:"code"`
	CreatedAt           time.Time `gorm:"column:created_at" json:"crated_at"`
	DepartaureTime      time.Time `gorm:"column:departaure_time" json:"departaure_time"`
	ArrivalTime         time.Time `gorm:"column:arrival_time" json:"arrival_time"`
	DepartaureAirportId int       `json:"-"`
	DepartaureAirport   Airport   `gorm:"foreignKey:DepartaureAirportId" json:"departaure_airport"`
	ArrivalAirportId    int       `json:"-"`
	ArrivalAirport      Airport   `gorm:"foreignKey:ArrivalAirportId" json:"arrival_airport"`
	Cost                float32   `gorm:"column:cost" json:"cost"`
}

// This interface is used so to have only one `ValidateFlight`
type FlightValidationInput interface {
	Departaure() (time.Time, int)
	Arrival() (time.Time, int)
}

// Struct used to get new data for a flight
type FlightInput struct {
	Code                string    `json:"code" binding:"required"`
	DepartaureTime      time.Time `json:"departaure_time" binding:"required"`
	ArrivalTime         time.Time `json:"arrival_time" binding:"required"`
	DepartaureAirportId int       `json:"departaure_airport_id" binding:"required"`
	ArrivalAirportId    int       `json:"arrival_airport_id" binding:"required"`
	Cost                float32   `json:"cost" binding:"required"`
}

func (in FlightInput) Departaure() (time.Time, int) {
	return in.DepartaureTime, in.DepartaureAirportId
}

func (in FlightInput) Arrival() (time.Time, int) {
	return in.ArrivalTime, in.ArrivalAirportId
}

// Struct used to get info on filter
type FlightFilterInput struct {
	DepartaureTime      time.Time `json:"departaure_time" binding:"required"`
	ArrivalTime         time.Time `json:"arrival_time" binding:"required"`
	DepartaureAirportId int       `json:"departaure_airport_id" binding:"required"`
	ArrivalAirportId    int       `json:"arrival_airport_id" binding:"required"`
}

func (in FlightFilterInput) Departaure() (time.Time, int) {
	return in.DepartaureTime, in.DepartaureAirportId
}

func (in FlightFilterInput) Arrival() (time.Time, int) {
	return in.ArrivalTime, in.ArrivalAirportId
}

// It validates data from `in` and returns a possible error or not
func ValidateFlight(db *gorm.DB, in FlightValidationInput) error {
	var departaure_airport Airport
	departaure_time, departaure_airport_id := in.Departaure()
	arrival_time, arrival_airport_id := in.Arrival()
	if err := db.Where("id = ?", departaure_airport_id).First(&departaure_airport).Error; err != nil {
		return errors.New("`departaure_airport_id` does not exist.")
	}

	var arrival_airport Airport
	if err := db.Where("id = ?", arrival_airport_id).First(&arrival_airport).Error; err != nil {
		return errors.New("`arrival_airport_id` does not exist.")
	}

	if departaure_airport_id == arrival_airport_id {
		return errors.New("`departaure_airport_id` can't be equals to `arrival_airport_id`")
	}

	if departaure_time.Equal(arrival_time) || departaure_time.After(arrival_time) {
		return errors.New("`departaure_time` can't be after or the same `arrival_time`")
	}

	return nil
}

// Returns a new Flight with the data from `in`. It should be called after
// `ValidateFlight(..., in)` method
func NewFlight(in FlightInput) Flight {
	return Flight{
		CreatedAt:           time.Now(),
		Code:                in.Code,
		DepartaureTime:      in.DepartaureTime,
		DepartaureAirportId: in.DepartaureAirportId,
		ArrivalTime:         in.ArrivalTime,
		ArrivalAirportId:    in.ArrivalAirportId,
		Cost:                in.Cost,
	}
}
