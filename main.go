package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
)

var (
	config = oauth2.Config{
		ClientID:     "<client id not yet configured>",
		ClientSecret: "<client secret not yet configured>",
		Scopes:       []string{"user.activity,user.metrics"},
		RedirectURL:  "https://withings.roessland.com/callback",
		// This points to our Authorization Server
		// if our Client ID and Client Secret are valid
		// it will attempt to authorize our user
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://account.withings.com/oauth2_user/authorize2",
			TokenURL: "https://account.withings.com/oauth2/token",
		},
	}

	token oauth2.Token
)


// Homepage
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Homepage Hit!")
	u := config.AuthCodeURL("xyfdsfdsz")
	http.Redirect(w, r, u, http.StatusFound)
}

// Authorize
func Authorize(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	state := r.Form.Get("state")
	if state != "xyfdsfdsz" {
		http.Error(w, "State invalid", http.StatusBadRequest)
		return
	}

	code := r.Form.Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshal token to JSON
	tokenJson, err := json.MarshalIndent(*token, "", " ")

	// Print token JSON to response
	_, err = w.Write(tokenJson)
	if err != nil {
		log.Println(err)
	}

	// Print token JSON to token cache
	f, err := os.Create("access_token.json")
	if err != nil {
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Println(err)
		}
	}
	_, err = f.Write(tokenJson)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	readConfig()

	// 1 - We attempt to hit our Homepage route
	// if we attempt to hit this unauthenticated, it
	// will automatically redirect to our Auth
	// server and prompt for login credentials
	http.HandleFunc("/", HomePage)

	// 2 - This displays our state, code and
	// token and expiry time that we get back
	// from our Authorization server
	http.HandleFunc("/callback", Authorize)

	// Use cached access token
	go func() {
		SleepGetSummary()
	}()

	// 3 - We start up our Client on port 9094
	log.Println("Client is running at 9094 port.")
	log.Fatal(http.ListenAndServe(":9094", nil))
}

func readConfig() {
	// Oauth app parameters
	config.ClientID = os.Getenv("WITHINGS_CLIENT_ID")
	config.ClientSecret = os.Getenv("WITHINGS_CLIENT_SECRET")

	// Read Oauth user access token into memory
	f, err := os.Open("access_token.json")
	if err != nil {
		log.Println(err)
	}
	err = json.NewDecoder(f).Decode(&token)
	if err != nil {
		log.Println(err)
	}
}