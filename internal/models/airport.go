package models

type Airport struct {
	Id        uint    `gorm:"column:id" json:"id"`
	Name      string  `gorm:"column:name" json:"name"`
	Location  string  `gorm:"column:location" json:"location"`
	Latitude  float32 `gorm:"column:latitude" json:"latitude"`
	Longitude float32 `gorm:"column:longitude" json:"longitude"`
}

type AirportInput struct {
	Name      string  `json:"name" binding:"required"`
	Location  string  `json:"location" binding:"required"`
	Latitude  float32 `json:"latitude" binding:"required"`
	Longitude float32 `json:"longitude" binding:"required"`
}

func NewAirport(in AirportInput) Airport {
	return Airport{
		Name:      in.Name,
		Location:  in.Location,
		Latitude:  in.Latitude,
		Longitude: in.Longitude,
	}
}
