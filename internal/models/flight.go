package models

import "time"

type Flight struct {
	Id                uint      `gorm:"column:id"`
	CreatedAt         time.Time `gorm:"column:created_at"`
	DepartaureTime    time.Time `gorm:"column:departaure_time"`
	ArrivalTime       time.Time `gorm:"column:arrival_time"`
	DepartaureAirport Airport
	ArrivalAirport    Airport
	Cost              float32 `form:"column:corst"`
}
