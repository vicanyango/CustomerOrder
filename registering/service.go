package registering

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/muhoro/log"
)

type RegisteringService interface {
	CreateCustomer(customer Customers) error
	CreateOrder(order Order) error
}

type RegisteringRepo interface {
	CreateCustomer(customer Customers) error
	ValidateCustomerIDAndGetPhoneNumber(customerID uuid.UUID) (phoneNumber string, err error)
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
	values := url.Values{}
	if err := s.validateOrder(order); err != nil {
		return err
	}
	uid, err := parseUUID(order.CustomerID)
	if err != nil {
		return err
	}
	phoneNumber, err := s.repo.ValidateCustomerIDAndGetPhoneNumber(uid)
	if err != nil {
		return err
	}
	if err := s.repo.CreateOrder(order, uid); err != nil {
		return err
	}
	msg := "order sent"
	username := "sandbox"
	values.Set("username", username)
	values.Set("message", msg)
	values.Set("to", phoneNumber)
	// values.Set("from", "CustomerOrder")
	atresponse, err := send(values)
	if !strings.Contains(strings.ToLower(atresponse.SMSMessageData.Message), "sent") {
		log.Error("Error encountered in sending SMS via AT: "+atresponse.SMSMessageData.Message, nil)
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

func send(data url.Values) (response, error) {
	var resppayload response

	url := "https://api.sandbox.africastalking.com/version1/messaging"
	apiKey := "371abf7ec701e55f5039a61805af9c668b5e27173b32741e947ec10117c0a87d"
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		log.Error("Failed to create HTTP request: "+err.Error(), nil)
		return resppayload, err
	}
	req.Header.Add("apiKey", apiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("HTTP request to microservice failed: "+err.Error(), req)
		return resppayload, err
	}

	defer resp.Body.Close()
	d, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 201 {
		return resppayload, errors.New(string(d))
	}

	err = xml.Unmarshal(d, &resppayload)

	return resppayload, nil
}
