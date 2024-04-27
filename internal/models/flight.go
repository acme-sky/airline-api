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

// Struct used to get info on filter. Filter fliths by airport code
type FlightFilterInput struct {
	DepartaureTime    time.Time `json:"departaure_time" binding:"required"`
	ArrivalTime       time.Time `json:"arrival_time" binding:"required"`
	DepartaureAirport string    `json:"departaure_airport" binding:"required"`
	ArrivalAirport    string    `json:"arrival_airport" binding:"required"`
}

func (in FlightFilterInput) Departaure() (time.Time, string) {
	return in.DepartaureTime, in.DepartaureAirport
}

func (in FlightFilterInput) Arrival() (time.Time, string) {
	return in.ArrivalTime, in.ArrivalAirport
}

// It validates data from `in` and returns a possible error or not
func ValidateFlight(db *gorm.DB, in FlightInput) error {
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

// It validates data from `in` and returns a possible error or not. If there's
// no error, returns an array of airport ids: first one for departaure airport
// and the latter for the arrival airport.
func ValidateFlightFilter(db *gorm.DB, in FlightFilterInput) (*[2]uint, error) {
	var departaure_airport Airport
	departaure_time, departaure_airport_code := in.Departaure()
	arrival_time, arrival_airport_code := in.Arrival()
	if err := db.Where("code = ?", departaure_airport_code).First(&departaure_airport).Error; err != nil {
		return nil, errors.New("`departaure_airport_id` does not exist.")
	}

	var arrival_airport Airport
	if err := db.Where("code = ?", arrival_airport_code).First(&arrival_airport).Error; err != nil {
		return nil, errors.New("`arrival_airport_id` does not exist.")
	}

	if departaure_airport.Id == arrival_airport.Id {
		return nil, errors.New("`departaure_airport_id` can't be equals to `arrival_airport_id`")
	}

	if departaure_time.Equal(arrival_time) || departaure_time.After(arrival_time) {
		return nil, errors.New("`departaure_time` can't be after or the same `arrival_time`")
	}

	airports := [2]uint{departaure_airport.Id, arrival_airport.Id}

	return &airports, nil
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
