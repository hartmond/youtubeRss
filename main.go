package main

import (
	"fmt"
	miniflux "miniflux.app/client"
)

const (
	youtubeSecretFile     = "youtube_secret.json"
	youtubeUserSecretFile = "youtube_user.json"
	minifluxSecretFile    = "miniflux_secret.json"
)

func main() {
	fmt.Println("starting...")
	minifluxClient, err := NewMinifluxClient(minifluxSecretFile)
	if err != nil {
		panic(err)
	}
	minifluxSubscriptions, err := minifluxClient.GetYoutubeSubscriptions()
	if err != nil {
		panic(err)
	}
	//fmt.Println(minifluxSubscriptions)

	youtubeClient, err := NewYoutubeClient(youtubeSecretFile, youtubeUserSecretFile)
	if err != nil {
		panic(err)
	}
	youtubeSubcriptions, err := youtubeClient.GetSubscriptions()
	if err != nil {
		panic(err)
	}

	// find newly subscripbed channels
	for _, elem := range youtubeSubcriptions {
		if !minifluxContains(elem, minifluxSubscriptions) {
			err = minifluxClient.Subscribe(elem)
			if err != nil {
				panic(err)
			}
			fmt.Println("subscribed: ", elem)
		}
	}

	// find unsubscribed channels
	for _, elem := range minifluxSubscriptions {
		if !youtubeContains(elem, youtubeSubcriptions) {
			err = minifluxClient.Unsubscribe(elem)
			if err != nil {
				panic(err)
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
