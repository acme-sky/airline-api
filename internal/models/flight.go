package models

import (
	"time"
)

type Flight struct {
	Id                  uint      `gorm:"column:id" json:"id"`
	CreatedAt           time.Time `gorm:"column:created_at" json:"crated_at"`
	DepartaureTime      time.Time `gorm:"column:departaure_time" json:"departaure_time"`
	ArrivalTime         time.Time `gorm:"column:arrival_time" json:"arrival_time"`
	DeparatureAirportId int
	DepartaureAirport   Airport `gorm:"foreignKey:DeparatureAirportId" json:"deparature_airport"`
	ArrivalAirportId    int
	ArrivalAirport      Airport `gorm:"foreignKey:ArrivalAirportId" json:"arrival_airport"`
	Cost                float32 `form:"column:cost" json:"cost"`
}
