package data

import (
	"CustomerOrder/registering"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewRepository(database *gorm.DB) CustomerRepository {
	return CustomerRepository{database}
}

func (repo CustomerRepository) CreateCustomer(customer registering.Customers) error {
	Customer := Customer{
		ID:          newUUID(),
		FirstName:   customer.FirstName,
		LastName:    customer.LastName,
		PhoneNumber: customer.PhoneNumber,
	}
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&Customer).Error; err != nil {
			return err
		}
		return nil
	})
}

func newUUID() uuid.UUID {
	uuid, _ := uuid.NewUUID()
	return uuid
}
