package main

import (
	"CustomerOrder/api"
	"CustomerOrder/data"
	"CustomerOrder/registering"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	oauthStateString  = "pseudo-random"
)

func main() {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=crm password=felixotieno sslmode=disable")

	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	services := initializeServices(db)
	route := mux.NewRouter()
	initializeRoutes(route, services)
	http.ListenAndServe(":2000", nil)
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
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleGoogleLogin)
	http.HandleFunc("/customerorder", handleGoogleCallback)
	route.HandleFunc("/customerorder/api/customer", api.CreateCustomer(s.registering)).Methods("POST")
	route.HandleFunc("/customerorder/api/order", api.CreateOrder(s.registering)).Methods("POST")
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:2000/customerorder",
		ClientID:     "http://198363440932-tk3225nqc33voqfntlsvmdhugp2fo9tc.apps.googleusercontent.com/",
		ClientSecret: "ytruxgpvPtqzVtaaBQWykxmj",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = `<html><body><a href="/login">Google Log In</a></body></html>`
	fmt.Fprintf(w, htmlIndex)
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	content, err := getUserInfo(r.FormValue("state"), r.FormValue("code"))
	if err != nil {
		fmt.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "Content: %s\n", content)
}

func getUserInfo(state string, code string) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, nil
}
