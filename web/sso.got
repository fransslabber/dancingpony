package rest_api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var ssogolang *oauth2.Config
var RandomString = "asdf"

func init() {

	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file: " + err.Error())
	}

	ssogolang = &oauth2.Config{
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{"email"}, // "https://www.googleapis.com/auth/userinfo.email"
		Endpoint:     google.Endpoint,

		// oauth2.Endpoint{  // google.Endpoint
		// 	AuthURL: os.Getenv("AUTH_URI"), // https://accounts.google.com/o/oauth2/auth
		// 	TokenURL: os.Getenv("TOKEN_URI"), // "https://accounts.google.com/o/oauth2/token"
		// },

	}
}

// create func SSO as http middleware
func SSO(w http.ResponseWriter, r *http.Request) {
	url := ssogolang.AuthCodeURL(RandomString)
	fmt.Println(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
