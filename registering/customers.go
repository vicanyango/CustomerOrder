package registering

import (
	"time"
)

type Customers struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
}

type Order struct {
	CustomerID string    `json:"customer_id"`
	Item       string    `json:"item"`
	Amount     float64   `json:"amount"`
	Time       time.Time `json:"time"`
}

type response struct {
	SMSMessageData struct {
		Message string `xml:"Message"`
	} `xml:"SMSMessageData"`
}
