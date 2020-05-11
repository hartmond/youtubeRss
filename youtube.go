package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

func youtubeTest() {
	service, err := getYoutubeClient()
	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	for pageToken := ""; ; {
		res, err := service.Subscriptions.List("snippet").Mine(true).MaxResults(10).PageToken(pageToken).Do()
		if err != nil {
			log.Fatalf("Error Listing YouTube Subscriptions: %v", err)
		}

		for _, elem := range res.Items {
			fmt.Println(elem.Snippet.ResourceId.ChannelId, elem.Snippet.Title)
		}

		pageToken = res.NextPageToken
		if pageToken == "" {
			break
		}
	}
}

func getYoutubeClient() (*youtube.Service, error) {
	scope := youtube.YoutubeReadonlyScope
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		return nil, fmt.Errorf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, scope)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse client secret file to config: %v", err)
	}

	f, err := os.Open("youtube-go.json")
	if err != nil {
		return nil, err
	}
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	httpClient := config.Client(ctx, token)

	return youtube.New(httpClient)
}
