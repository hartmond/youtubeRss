package main

import (
	"fmt"
)

func main() {
	fmt.Println("starting...")
	
	minifluxClient, err := NewMinifluxClient("https://xxx/", "xxx")
	if err != nil {
		panic(err)
	}
	minifluxSubscriptions, err := minifluxClient.GetYoutubeSubscriptions()
	if err != nil {
		panic(err)
	}
	fmt.Println(minifluxSubscriptions)
	
	youtubeClient, err := NewYoutubeClient()
	if err != nil {
		panic(err)
	}
	youtubeSubcriptions, err := youtubeClient.GetSubscriptions()
	if err != nil {
		panic(err)
	}
	fmt.Println(youtubeSubcriptions)
}
