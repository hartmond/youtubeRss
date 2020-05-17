package main

import (
	"fmt"
	"time"

	miniflux "miniflux.app/client"
)

const (
	updateInterval        = time.Hour
	youtubeSecretFile     = "youtube_secret.json"
	youtubeUserSecretFile = "youtube_user.json"
	minifluxSecretFile    = "miniflux_secret.json"
)

var (
	minifluxClient *MinifluxClient
	youtubeClient  *YoutubeClient
)

func main() {
	fmt.Println("loading configurations and initializing clients")

	var err error
	minifluxClient, err = NewMinifluxClient(minifluxSecretFile)
	if err != nil {
		fmt.Println("Error initializing Miniflux Client: %v", err)
		return
	}
	youtubeClient, err = NewYoutubeClient(youtubeSecretFile, youtubeUserSecretFile)
	if err != nil {
		fmt.Println("Error initializing YouTube Client: %v", err)
		return
	}

	fmt.Println("staring main loop")
	ticker := time.NewTicker(updateInterval)

	for {
		fmt.Println("starting update procedure")
		updateFeeds()

		<-ticker.C // wait until next update
	}
}

func updateFeeds() {
	// get current subscriptions
	minifluxSubscriptions, err := minifluxClient.GetYoutubeSubscriptions()
	if err != nil {
		fmt.Println("Error receiving current subscriptions from Miniflux: %v", err)
		return
	}

	youtubeSubcriptions, err := youtubeClient.GetSubscriptions()
	if err != nil {
		fmt.Println("Error receiving current subscriptions from YouTube: %v", err)
		return
	}

	// find newly subscripbed channels
	for _, elem := range youtubeSubcriptions {
		if !minifluxContains(elem, minifluxSubscriptions) {
			err = minifluxClient.Subscribe(elem)
			if err != nil {
				fmt.Println("Error on subscribe of %v: %v", elem, err)
				continue
			}
			fmt.Println("subscribed: ", elem)
		}
	}

	// find unsubscribed channels
	for _, elem := range minifluxSubscriptions {
		if !youtubeContains(elem, youtubeSubcriptions) {
			err = minifluxClient.Unsubscribe(elem)
			if err != nil {
				fmt.Println("Error on subscribe of %v: %v", elem.FeedURL, err)
				continue
			}
			fmt.Println("unsubscripbed: ", elem.FeedURL)
		}
	}
}

func minifluxContains(val string, list []miniflux.Feed) bool {
	for _, elem := range list {
		if val == elem.FeedURL {
			return true
		}
	}
	return false
}

func youtubeContains(val miniflux.Feed, list []string) bool {
	for _, elem := range list {
		if val.FeedURL == elem {
			return true
		}
	}
	return false
}
