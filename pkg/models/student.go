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

// Hooks
func (s *Student) BeforeCreate(_ *gorm.DB) (err error) {
	if !s.IsValid() {
		err = errors.New("student's data is not valid")
	}
	return
}

// CRUD
func CreateStudent(db *gorm.DB, s *Student) error {
	result := db.Create(s)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func FindFirstStudent(db *gorm.DB, s *Student, name string) error {
	result := db.First(s, "firstname = ?", name)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetAllStudent(db *gorm.DB, s *[]Student) error {
	result := db.Find(s)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func SaveStudent(db *gorm.DB, s *Student) error {
	////
	//var tempStudent Student
	//findError := FindFirstStudent(db, &tempStudent, name)

	result := db.Save(s)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
