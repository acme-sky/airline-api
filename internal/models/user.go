package models

import "gorm.io/gorm"

// User model
// We ignore all the implementation for users having a manually creation. This
// model is used only for login.
type User struct {
	gorm.Model
	Username string `gorm:"column:username" gorm:"uniqueIndex"`
	Password string `gorm:"column:password"`
}
