package api

import (
	"CustomerOrder/registering"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/muhoro/log"
)

func CreateCustomer(service registering.RegisteringService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		customer := registering.Customers{}
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error(err.Error(), nil)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = json.Unmarshal(requestBody, &customer); err != nil {
			log.Error(err.Error(), &customer)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err = service.CreateCustomer(customer); err != nil {
			log.Error(err.Error(), &customer)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
