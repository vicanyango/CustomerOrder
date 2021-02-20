package registering

import (
	"errors"

	"github.com/google/uuid"
)

type RegisteringService interface {
	CreateCustomer(customer Customers) error
	CreateOrder(order Order) error
}

type RegisteringRepo interface {
	CreateCustomer(customer Customers) error
	ValidateCustomerID(customerID uuid.UUID) error
	CreateOrder(order Order, custID uuid.UUID) error
}

type service struct {
	repo RegisteringRepo
}

func NewRegisteringService(repo RegisteringRepo) RegisteringService {
	return &service{repo: repo}
}

func (s *service) CreateCustomer(customers Customers) error {
	err := s.validateCustomerInfo(customers)
	if err != nil {
		return err
	}
	erro := s.repo.CreateCustomer(customers)
	if erro != nil {
		return erro
	}
	return nil

}

func (s *service) validateCustomerInfo(cust Customers) error {
	if len(cust.FirstName) == 0 {
		return errors.New("first name cannot be empty")
	}
	if len(cust.LastName) == 0 {
		return errors.New("last name cannot be empty")
	}
	if len(cust.PhoneNumber) == 0 {
		return errors.New("phone number cannot be empty")
	}
	return nil
}

func (s *service) CreateOrder(order Order) error {
	if err := s.validateOrder(order); err != nil {
		return err
	}
	uid, err := parseUUID(order.CustomerID)
	if err != nil {
		return err
	}
	if err := s.repo.ValidateCustomerID(uid); err != nil {
		return err
	}
	if err := s.repo.CreateOrder(order, uid); err != nil {
		return err
	}
	return nil
}

func (s *service) validateOrder(order Order) error {
	if order.Amount == 0 {
		return errors.New("item amount cannot be empty")
	}
	if len(order.Item) == 0 {
		return errors.New("item caanot be empty")
	}
	if len(order.CustomerID) == 0 {
		return errors.New("customer id cannot be nil")
	}
	return nil
}

func parseUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}
