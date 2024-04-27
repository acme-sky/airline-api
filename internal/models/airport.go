package models

// Airport model
type Airport struct {
	Id        uint    `gorm:"column:id" json:"id"`
	Name      string  `gorm:"column:name" json:"name"`
	Code      string  `gorm:"column:code" json:"code"`
	Location  string  `gorm:"column:location" json:"location"`
	Latitude  float32 `gorm:"column:latitude" json:"latitude"`
	Longitude float32 `gorm:"column:longitude" json:"longitude"`
}

// Struct used to get new data for an airport
type AirportInput struct {
	Name      string  `json:"name" binding:"required"`
	Code      string  `json:"code" binding:"required"`
	Location  string  `json:"location" binding:"required"`
	Latitude  float32 `json:"latitude" binding:"required"`
	Longitude float32 `json:"longitude" binding:"required"`
}

// Returns a new Airport with the data from `in`
func NewAirport(in AirportInput) Airport {
	return Airport{
		Name:      in.Name,
		Code:      in.Code,
		Location:  in.Location,
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
	}
}
