package main

import (
	"CustomerOrder/api"
	"CustomerOrder/data"
	"CustomerOrder/registering"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/muhoro/log"
)

// var Db *gorm.DB

func main() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=crm password=felixotieno sslmode=disable")

	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	services := initializeServices(db)
	route := mux.NewRouter()
	initializeRoutes(route, services)
	log.Fatal(http.ListenAndServe(":2000", route).Error(), route)
}

type services struct {
	registering registering.RegisteringService
}

func initializeServices(db *gorm.DB) services {
	dbrepo := data.NewRepository(db)
	s := services{}
	s.registering = registering.NewRegisteringService(dbrepo)

	return s
}

func initializeRoutes(route *mux.Router, s services) {
	route.HandleFunc("/customerorder/api/customer", api.CreateCustomer(s.registering)).Methods("POST")
	route.HandleFunc("/customerorder/api/order", api.CreateOrder(s.registering)).Methods("POST")
}
