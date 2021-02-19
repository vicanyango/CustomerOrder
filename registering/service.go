package registering

import "errors"

type RegisteringService interface {
	CreateCustomer(customer Customers) error
}

type RegisteringRepo interface {
	CreateCustomer(customer Customers) error
}

type service struct {
	repo RegisteringRepo
}

func NewRegisteringService(repo RegisteringRepo) RegisteringService {
	return &service{repo: repo}
}

func (s *service) CreateCustomer(customers Customers) error {
	err := validateCustomerInfo(customers)
	if err != nil {
		return err
	}
	erro := s.repo.CreateCustomer(customers)
	if erro != nil {
		return erro
	}
	return nil

}

func validateCustomerInfo(cust Customers) error {
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
