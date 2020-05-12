package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

type YoutubeClient struct {
	client *youtube.Service
}

func NewYoutubeClient() (*YoutubeClient, error) {
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
	youtubeClient, err := youtube.New(httpClient)
	if err != nil {
		return nil, err
	}
	return &YoutubeClient{youtubeClient}, nil
}

func (client *YoutubeClient) GetSubscriptions() ([]string, error) {
	results := []string{}

	for pageToken := ""; ; {
		res, err := client.client.Subscriptions.List("snippet").Mine(true).MaxResults(10).PageToken(pageToken).Do()
		if err != nil {
			return nil, err
		}

		for _, elem := range res.Items {
			//fmt.Println(elem.Snippet.ResourceId.ChannelId, elem.Snippet.Title)
			results = append(results, fmt.Sprintf("https://www.youtube.com/feeds/videos.xml?channel_id=%s", elem.Snippet.ResourceId.ChannelId))
		}

		pageToken = res.NextPageToken
		if pageToken == "" {
			return results, nil
		}
	}
}
