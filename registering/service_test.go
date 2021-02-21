package registering

import (
	"errors"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockRegisteringRepo struct {
	mock.Mock
}

func (m *mockRegisteringRepo) CreateCustomer(customer Customers) error {
	args := m.Called(customer)
	return args.Error(0)
}
func (m *mockRegisteringRepo) ValidateCustomerIDAndGetPhoneNumber(customerID uuid.UUID) (string, error) {
	args := m.Called(customerID)
	return args.String(0), args.Error(1)
}
func (m *mockRegisteringRepo) CreateOrder(order Order, custID uuid.UUID) error {
	args := m.Called(order, custID)
	return args.Error(0)
}

type RegisteringTestSuite struct {
	suite.Suite
	service RegisteringService
	repo    *mockRegisteringRepo
}

func (s *RegisteringTestSuite) SetupTest() {
	// create mocks and initialize service for testing
	s.repo = new(mockRegisteringRepo)
	s.service = NewRegisteringService(s.repo)
}

// TearDownTest ensures that all expected tests were run
func (s *RegisteringTestSuite) TearDownTest() {
	s.repo.AssertExpectations(s.T())
}

func (s *RegisteringTestSuite) TestCreateCustomerSuccess() {
	cust := fakeCustomerData()
	s.repo.On("CreateCustomer", mock.Anything).Return(nil)
	err := s.service.CreateCustomer(cust)
	s.Nil(err)
	s.repo.AssertCalled(s.T(), "CreateCustomer", cust)
}

func (s *RegisteringTestSuite) TestCreateCustomerRepoFailure() {
	cust := fakeCustomerData()
	s.repo.On("CreateCustomer", mock.Anything).Return(errors.New("Failed to saved customer in the database"))
	err := s.service.CreateCustomer(cust)
	s.NotNil(err)
	s.True(strings.Contains(err.Error(), string("Failed to saved customer")))
}

func (s *RegisteringTestSuite) TestValidateCustomerInfoWrongFirstNameFailure() {
	cust := fakeCustomerData()
	cust.FirstName = ""
	err := s.service.CreateCustomer(cust)
	s.NotNil(err)
	s.True(strings.Contains(err.Error(), string("first name cannot be empty")))
}

func (s *RegisteringTestSuite) TestValidateCustomerInfoWrongLastNameFailure() {
	cust := fakeCustomerData()
	cust.LastName = ""
	err := s.service.CreateCustomer(cust)
	s.NotNil(err)
	s.True(strings.Contains(err.Error(), string("last name cannot be empty")))
}

func (s *RegisteringTestSuite) TestValidateCustomerInfoWrongPhoneNumberFailure() {
	cust := fakeCustomerData()
	cust.PhoneNumber = ""
	err := s.service.CreateCustomer(cust)
	s.NotNil(err)
	s.True(strings.Contains(err.Error(), string("phone number cannot be empty")))
}

func (s *RegisteringTestSuite) TestCreateOrderSuccess() {
	orderInfo := fakeCustomerOrderInfo()
	s.repo.On("ValidateCustomerIDAndGetPhoneNumber", mock.Anything).Return("0715893271", nil)
	s.repo.On("CreateOrder", mock.Anything, mock.Anything).Return(nil)
	err := s.service.CreateOrder(orderInfo)
	s.Nil(err)
}

func (s *RegisteringTestSuite) TestValidateCustomerIDRepoFailure() {
	orderInfo := fakeCustomerOrderInfo()
	s.repo.On("ValidateCustomerIDAndGetPhoneNumber", mock.Anything).Return("", errors.New("Failed to get customer Id in the database"))
	err := s.service.CreateOrder(orderInfo)
	s.NotNil(err)
	s.True(strings.Contains(err.Error(), string("Failed to get customer Id")))
}

func (s *RegisteringTestSuite) TestCreateOrderRepoFailure() {
	orderInfo := fakeCustomerOrderInfo()
	s.repo.On("ValidateCustomerIDAndGetPhoneNumber", mock.Anything).Return("0715893271", nil)
	s.repo.On("CreateOrder", mock.Anything, mock.Anything).Return(errors.New("Failed to create customer order in the database"))
	err := s.service.CreateOrder(orderInfo)
	s.NotNil(err)
	s.True(strings.Contains(err.Error(), string("Failed to create customer order")))
}

func (s *RegisteringTestSuite) TestValidateOrderAmountFaliure() {
	orderInfo := fakeCustomerOrderInfo()
	orderInfo.Amount = 0
	err := s.service.CreateOrder(orderInfo)
	s.NotNil(err)
	s.True(strings.Contains(err.Error(), string("item amount cannot be empty")))
}

func (s *RegisteringTestSuite) TestValidateOrderItemFaliure() {
	orderInfo := fakeCustomerOrderInfo()
	orderInfo.Item = ""
	err := s.service.CreateOrder(orderInfo)
	s.NotNil(err)
	s.True(strings.Contains(err.Error(), string("item cannot be empty")))
}

func (s *RegisteringTestSuite) TestValidateCustomerIDFaliure() {
	orderInfo := fakeCustomerOrderInfo()
	orderInfo.CustomerID = ""
	err := s.service.CreateOrder(orderInfo)
	s.NotNil(err)
	s.True(strings.Contains(err.Error(), string("customer id cannot be nil")))
}

func fakeCustomerData() Customers {
	return Customers{
		FirstName:   "Test",
		LastName:    "Customer",
		PhoneNumber: "0715893271",
	}
}

func fakeCustomerOrderInfo() Order {
	return Order{
		CustomerID: "281e9818-34e8-49ec-b6ee-c282c38202be",
		Item:       "Milk",
		Amount:     50,
	}
}

func TestRegisteringService(t *testing.T) {
	suite.Run(t, new(RegisteringTestSuite))
}
