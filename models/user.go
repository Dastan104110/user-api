package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name string `json:"name" gorm:"unique;not null"`
	Age  int    `json:"age"`
}
