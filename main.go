package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mrjones/oauth"
)

// This HTTP app authenticates you against an app on Twitter and gives
// back your access credentials.

var creds = oauth.NewConsumer(
	os.Getenv("TWITTER_CONSUMER_KEY"),
	os.Getenv("TWITTER_CONSUMER_SECRET"),
	oauth.ServiceProvider{
		RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
		AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
		AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
	},
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "")
	})
	http.HandleFunc("/auth", startauth)
	http.HandleFunc("/callback", callback)

	host := ":" + os.Getenv("PORT")

	log.Printf("Listening on %s...", host)
	log.Fatal(http.ListenAndServe(host, nil))
}

var tokens = make(map[string]*oauth.RequestToken)

func startauth(w http.ResponseWriter, r *http.Request) {
	reqToken, url, err := creds.GetRequestTokenAndUrl("")
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	tokens[reqToken.Token] = reqToken
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func callback(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	verifier := r.FormValue("oauth_verifier")
	token := r.FormValue("oauth_token")

	pending, ok := tokens[token]
	if !ok {
		fmt.Fprintf(w, "No pending authorization found")
		return
	}

	delete(tokens, token)

	accessToken, err := creds.AuthorizeToken(pending, verifier)

	fmt.Fprintf(w, "%+v %v", accessToken, err)
}
