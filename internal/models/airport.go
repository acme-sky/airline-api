package models

type Airport struct {
	Id        uint    `gorm:"column:id"`
	Name      string  `gorm:"column:name"`
	Location  string  `gorm:"column:location"`
	Latitude  float32 `gorm:"column:latitude"`
	Longitude float32 `gorm:"column:longitude"`
}
