package data

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID          uuid.UUID `gorm:"column:id"`
	FirstName   string    `gorm:"column:first_name"`
	LastName    string    `gorm:"column:last_name"`
	PhoneNumber string    `gorm:"column:phone_number"`
}

type Order struct {
	ID         uuid.UUID `gorm:"column:id"`
	CustomerID uuid.UUID `gorm:"column:customer_id"`
	Item       string    `gorm:"column:item"`
	Amount     float64   `gorm:"column:amount"`
	Time       time.Time `gorm:"column:order_time"`
}
