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
}
