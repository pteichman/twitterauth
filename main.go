package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mrjones/oauth"
	"golang.org/x/oauth2"
)

// This HTTP app authenticates you against an OAuth app and gives back
// your access credentials.

var twittercreds = oauth.NewConsumer(
	os.Getenv("TWITTER_CONSUMER_KEY"),
	os.Getenv("TWITTER_CONSUMER_SECRET"),
	oauth.ServiceProvider{
		RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
		AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
		AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
	},
)

var instaendpoint = oauth2.Endpoint{
	AuthURL:  "https://instagram.com/oauth/authorize",
	TokenURL: "https://api.instagram.com/oauth/access_token",
}

var instaconfig = oauth2.Config{
	os.Getenv("INSTAGRAM_CONSUMER_KEY"),
	os.Getenv("INSTAGRAM_CONSUMER_SECRET"),
	instaendpoint,
	"https://still-citadel-7423.herokuapp.com/instagram/callback",
	[]string{"basic", "likes", "public_content"},
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "")
	})
	http.HandleFunc("/twitter/auth", twitterstartauth)
	http.HandleFunc("/twitter/callback", twittercallback)
	http.HandleFunc("/instagram/auth", instagramstartauth)
	http.HandleFunc("/instagram/callback", instagramcallback)

	host := ":" + os.Getenv("PORT")

	log.Printf("Listening on %s...", host)
	log.Fatal(http.ListenAndServe(host, nil))
}

var tokens = make(map[string]*oauth.RequestToken)

func twitterstartauth(w http.ResponseWriter, r *http.Request) {
	reqToken, url, err := twittercreds.GetRequestTokenAndUrl("")
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	tokens[reqToken.Token] = reqToken
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func twittercallback(w http.ResponseWriter, r *http.Request) {
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

	accessToken, err := twittercreds.AuthorizeToken(pending, verifier)

	fmt.Fprintf(w, "%+v %v", accessToken, err)
}

func instagramstartauth(w http.ResponseWriter, r *http.Request) {
	state := "asdfasdf"
	url := instaconfig.AuthCodeURL(state, oauth2.AccessTypeOffline)

	log.Println(url)

	//	tokens[state] = state
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func instagramcallback(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	code := r.FormValue("code")

	// pending, ok := tokens[token]
	// if !ok {
	// 	fmt.Fprintf(w, "No pending authorization found")
	// 	return
	// }

	//	delete(tokens, token)

	accessToken, err := instaconfig.Exchange(oauth2.NoContext, code)

	fmt.Fprintf(w, "%+v %v", accessToken, err)
}
