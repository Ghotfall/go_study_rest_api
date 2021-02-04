package models

import (
	"errors"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Firstname string
	Lastname  string
	Group     string
	Course    uint
}

func (s *Student) IsValid() bool {
	return s.Firstname != "" && s.Lastname != "" && s.Group != "" && s.Course != 0
}

func (s *Student) BeforeCreate(_ *gorm.DB) (err error) {
	if !s.IsValid() {
		err = errors.New("student's data is not valid")
	}
	return
}
