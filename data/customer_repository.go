package data

import (
	"CustomerOrder/registering"
	"time"

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

func (repo CustomerRepository) ValidateCustomerIDAndGetPhoneNumber(customerID uuid.UUID) (string, error) {
	customer := Customer{}
	err := repo.db.Debug().Where("id=?", customerID).Find(&customer).Error
	if err != nil {
		return "", err
	}
	return customer.PhoneNumber, nil
}

func (repo CustomerRepository) CreateOrder(order registering.Order, custID uuid.UUID) error {
	Order := Order{
		ID:         newUUID(),
		CustomerID: custID,
		Item:       order.Item,
		Amount:     order.Amount,
		Time:       time.Now().UTC(),
	}
	return repo.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&Order).Error; err != nil {
			return err
		}
		return nil
	})
}

func newUUID() uuid.UUID {
	uuid, _ := uuid.NewUUID()
	return uuid
}
