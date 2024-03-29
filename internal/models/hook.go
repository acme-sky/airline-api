package models

import (
	"time"
)

// Hook model
type Hook struct {
	Id        uint      `gorm:"column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"crated_at"`
	Name      string    `gorm:"column:name" json:"name"`
	Endpoint  string    `gorm:"column:endpoint" json:"endpoint"`
}

// Struct used to get new data for a hook
type HookInput struct {
	Name     string `json:"name" binding:"required"`
	Endpoint string `json:"endpoint" binding:"required"`
}

// Struct used to send a new request for a selected flight
type OffertInput struct {
	FlightId int `json:"flight_id" binding:"required"`
}

// Returns a new Hook with the data from `in`
func NewHook(in HookInput) Hook {
	return Hook{
		Name:      in.Name,
		CreatedAt: time.Now(),
		Endpoint:  in.Endpoint,
	}
}
