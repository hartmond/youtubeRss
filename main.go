package main

import (
	"log"
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
	log.Println("loading configurations and initializing clients")

	var err error
	minifluxClient, err = NewMinifluxClient(minifluxSecretFile)
	if err != nil {
		log.Printf("Error initializing Miniflux Client: %v", err)
		return
	}
	youtubeClient, err = NewYoutubeClient(youtubeSecretFile, youtubeUserSecretFile)
	if err != nil {
		log.Printf("Error initializing YouTube Client: %v", err)
		return
	}

	log.Println("staring main loop")
	ticker := time.NewTicker(updateInterval)

	for {
		log.Println("starting update procedure")
		updateFeeds()

		<-ticker.C // wait until next update
	}
}

func updateFeeds() {
	// get current subscriptions
	minifluxSubscriptions, err := minifluxClient.GetYoutubeSubscriptions()
	if err != nil {
		log.Printf("Error receiving current subscriptions from Miniflux: %v", err)
		return
	}

	youtubeSubscriptions, err := youtubeClient.GetSubscriptions()
	if err != nil {
		log.Printf("Error receiving current subscriptions from YouTube: %v", err)
		return
	}

	// find newly subscripbed channels
	for _, elem := range youtubeSubscriptions {
		if !minifluxContains(elem, minifluxSubscriptions) {
			err = minifluxClient.Subscribe(elem)
			if err != nil {
				log.Printf("Error on subscribe of %v: %v", elem, err)
				continue
			}
			log.Println("subscribed: ", elem)
		}
	}

	// find unsubscribed channels
	for _, elem := range minifluxSubscriptions {
		if !youtubeContains(elem, youtubeSubscriptions) {
			err = minifluxClient.Unsubscribe(elem)
			if err != nil {
				log.Printf("Error on unsubscribe of %v: %v", elem.FeedURL, err)
				continue
			}
			log.Println("unsubscripbed: ", elem.FeedURL)
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

