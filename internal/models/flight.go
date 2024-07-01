package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Flight model
type Flight struct {
	Id                 uint      `gorm:"column:id" json:"id"`
	Code               string    `gorm:"code" json:"code"`
	CreatedAt          time.Time `gorm:"column:created_at" json:"crated_at"`
	DepartureTime      time.Time `gorm:"column:departure_time" json:"departure_time"`
	ArrivalTime        time.Time `gorm:"column:arrival_time" json:"arrival_time"`
	DepartureAirportId int       `json:"-"`
	DepartureAirport   Airport   `gorm:"foreignKey:DepartureAirportId" json:"departure_airport"`
	ArrivalAirportId   int       `json:"-"`
	ArrivalAirport     Airport   `gorm:"foreignKey:ArrivalAirportId" json:"arrival_airport"`
	Cost               float32   `gorm:"column:cost" json:"cost"`
}

// Struct used to get new data for a flight
type FlightInput struct {
	Code               string    `json:"code" binding:"required"`
	DepartureTime      time.Time `json:"departure_time" binding:"required"`
	ArrivalTime        time.Time `json:"arrival_time" binding:"required"`
	DepartureAirportId int       `json:"departure_airport_id" binding:"required"`
	ArrivalAirportId   int       `json:"arrival_airport_id" binding:"required"`
	Cost               float32   `json:"cost" binding:"required"`
}

func (in FlightInput) Departure() (time.Time, int) {
	return in.DepartureTime, in.DepartureAirportId
}

func (in FlightInput) Arrival() (time.Time, int) {
	return in.ArrivalTime, in.ArrivalAirportId
}

// Struct used to get info on filter. Filter fliths by airport code
type FlightFilterInput struct {
	Code             *string   `json:"code"`
	DepartureTime    time.Time `json:"departure_time" binding:"required"`
	ArrivalTime      time.Time `json:"arrival_time" binding:"required"`
	DepartureAirport string    `json:"departure_airport" binding:"required"`
	ArrivalAirport   string    `json:"arrival_airport" binding:"required"`
}

func (in FlightFilterInput) Departure() (time.Time, string) {
	return in.DepartureTime, in.DepartureAirport
}

func (in FlightFilterInput) Arrival() (time.Time, string) {
	return in.ArrivalTime, in.ArrivalAirport
}

// It validates data from `in` and returns a possible error or not
func ValidateFlight(db *gorm.DB, in FlightInput) error {
	var departure_airport Airport
	departure_time, departure_airport_id := in.Departure()
	arrival_time, arrival_airport_id := in.Arrival()
	if err := db.Where("id = ?", departure_airport_id).First(&departure_airport).Error; err != nil {
		return errors.New("`departure_airport_id` does not exist.")
	}

	var arrival_airport Airport
	if err := db.Where("id = ?", arrival_airport_id).First(&arrival_airport).Error; err != nil {
		return errors.New("`arrival_airport_id` does not exist.")
	}

	if departure_airport_id == arrival_airport_id {
		return errors.New("`departure_airport_id` can't be equals to `arrival_airport_id`")
	}

	if departure_time.Equal(arrival_time) || departure_time.After(arrival_time) {
		return errors.New("`departure_time` can't be after or the same `arrival_time`")
	}

	return nil
}

// It validates data from `in` and returns a possible error or not. If there's
// no error, returns an array of airport ids: first one for departure airport
// and the latter for the arrival airport.
func ValidateFlightFilter(db *gorm.DB, in FlightFilterInput) (*[2]uint, error) {
	var departure_airport Airport
	_, departure_airport_code := in.Departure()
	_, arrival_airport_code := in.Arrival()
	if err := db.Where("code = ?", departure_airport_code).First(&departure_airport).Error; err != nil {
		return nil, errors.New("`departure_airport_id` does not exist.")
	}

	var arrival_airport Airport
	if err := db.Where("code = ?", arrival_airport_code).First(&arrival_airport).Error; err != nil {
		return nil, errors.New("`arrival_airport_id` does not exist.")
	}

	if departure_airport.Id == arrival_airport.Id {
		return nil, errors.New("`departure_airport_id` can't be equals to `arrival_airport_id`")
	}

	airports := [2]uint{departure_airport.Id, arrival_airport.Id}

	return &airports, nil
}

// Returns a new Flight with the data from `in`. It should be called after
// `ValidateFlight(..., in)` method
func NewFlight(in FlightInput) Flight {
	return Flight{
		CreatedAt:          time.Now(),
		Code:               in.Code,
		DepartureTime:      in.DepartureTime,
		DepartureAirportId: in.DepartureAirportId,
		ArrivalTime:        in.ArrivalTime,
		ArrivalAirportId:   in.ArrivalAirportId,
		Cost:               in.Cost,
	}
}
