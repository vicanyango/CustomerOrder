package data

import "github.com/google/uuid"

type Customer struct {
	ID          uuid.UUID `gorm:"column:id"`
	FirstName   string    `gorm:"column:first_name"`
	LastName    string    `gorm:"column:last_name"`
	PhoneNumber string    `gorm:column:phone_number"`
}
