package main

import (
	"CustomerOrder/api"
	"CustomerOrder/data"
	"CustomerOrder/registering"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	// "io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}

func main() {
	connectionString := getEnv("connection_string")
	db, err := gorm.Open("postgres", connectionString)

	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	services := initializeServices(db)
	route := mux.NewRouter()
	initializeRoutes(route, services)
	port := os.Getenv("PORT")
	if port == "" {
		port = "2000"
	}
	http.ListenAndServe(":"+port, nil)
}

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
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

	// Root route
	// Simply returns a link to the login route
	http.HandleFunc("/", rootHandler)

	// Login route
	http.HandleFunc("/login/github/", githubLoginHandler)

	// Github callback
	http.HandleFunc("/customerorder/api/customer", githubCallbackHandler)

	// Route where the authenticated user is redirected to
	http.HandleFunc("/loggedin", func(w http.ResponseWriter, r *http.Request) {
		loggedinHandler(w, r, "")
	})
	http.HandleFunc("/customerorder/api/customers", api.CreateCustomer(s.registering))
	http.HandleFunc("/customerorder/api/order", api.CreateOrder(s.registering))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<a href="https://github.com/login/oauth/authorize?client_id=84bee5d27ad1a3270adb&redirect_uri=http://localhost:2000/customerorder/api/customer">LOGIN</a>`)
}

func loggedinHandler(w http.ResponseWriter, r *http.Request, githubData string) {
	if githubData == "" {
		// Unauthorized users get an unauthorized message
		fmt.Fprintf(w, "UNAUTHORIZED!")
		return
	}

	w.Header().Set("Content-type", "application/json")

	// Prettifying the json
	var prettyJSON bytes.Buffer
	// json.indent is a library utility function to prettify JSON indentation
	parserr := json.Indent(&prettyJSON, []byte(githubData), "", "\t")
	if parserr != nil {
		log.Panic("JSON parse error")
	}

	// Return the prettified JSON as a string
	fmt.Fprintf(w, string(prettyJSON.Bytes()))
}

func githubLoginHandler(w http.ResponseWriter, r *http.Request) {
	githubClientID := getEnv("CLIENT_ID")
	redirectURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s", githubClientID, "http://localhost:2000/customerorder/api/customer")

	http.Redirect(w, r, redirectURL, 301)
}

func githubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	githubAccessToken := getGithubAccessToken(code)

	githubData := getGithubData(githubAccessToken)

	loggedinHandler(w, r, githubData)
}

func getGithubData(accessToken string) string {
	req, reqerr := http.NewRequest("GET", "https://api.github.com/user", nil)
	if reqerr != nil {
		log.Panic("API Request creation failed")
	}

	authorizationHeaderValue := fmt.Sprintf("token %s", accessToken)
	req.Header.Set("Authorization", authorizationHeaderValue)

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	respbody, _ := ioutil.ReadAll(resp.Body)

	return string(respbody)
}

func getGithubAccessToken(code string) string {

	clientID := getEnv("CLIENT_ID")
	clientSecret := getEnv("CLIENT_SECRET")

	requestBodyMap := map[string]string{"client_id": clientID, "client_secret": clientSecret, "code": code}
	requestJSON, _ := json.Marshal(requestBodyMap)

	req, reqerr := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewBuffer(requestJSON))
	if reqerr != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, resperr := http.DefaultClient.Do(req)
	if resperr != nil {
		log.Panic("Request failed")
	}

	respbody, _ := ioutil.ReadAll(resp.Body)

	// Represents the response received from Github
	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	var ghresp githubAccessTokenResponse
	json.Unmarshal(respbody, &ghresp)

	return ghresp.AccessToken
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatal("key not defined in .env file")
	}
	return value
}
