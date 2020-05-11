package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/youtube/v3"
)

func main() {
	scope := youtube.YoutubeReadonlyScope

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, scope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	config.RedirectURL = "urn:ietf:wg:oauth:2.0:oob"

	cacheFile := "youtube-go.json"

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	var code string
	fmt.Printf("Open the following link in your browser:\n\n%v\n\n", authURL)
	fmt.Println("Enter the displayed token:")

	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	fmt.Println(authURL)

	token, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token %v", err)
	}

	if err == nil {
		fmt.Println("Saving token")
		fmt.Printf("Saving credential file to: %s\n", cacheFile)
		f, err := os.OpenFile(cacheFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			log.Fatalf("Unable to cache oauth token: %v", err)
		}
		defer f.Close()
		json.NewEncoder(f).Encode(token)
	}

}
