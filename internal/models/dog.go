package models

import (
	"time"

	"github.com/google/uuid"
)

type Dog struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" validate:"required,min=2,max=24" gorm:"check:name <> ''"`
	Breed     string    `json:"breed" validate:"required,min=2,max=24" gorm:"check:breed <> ''"`
	Color     string    `json:"color" validate:"required,min=2,max=24" gorm:"check:color <> ''"`
	Birthdate time.Time `json:"birthdate" validate:"required"`
	Weight    int       `json:"weight" validate:"required,gte=1,lt=300" gorm:"check:weight > 0"`
}
