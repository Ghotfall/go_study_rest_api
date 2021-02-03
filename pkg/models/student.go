package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	Firstname string
	Lastname  string
}
