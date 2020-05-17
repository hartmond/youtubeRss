package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"

	miniflux "miniflux.app/client"
)

type minifluxConfig struct {
	Url string
	Token string
}

type MinifluxClient struct {
	client          *miniflux.Client
	youtubeCategory *miniflux.Category
}

func NewMinifluxClient(minifluxSecretFile string) (*MinifluxClient, error) {
	// parse config
	minifluxSecretFileHandle, err := os.Open(minifluxSecretFile)
	if err != nil {
		return nil, fmt.Errorf("Miniflux Config file could not be opened: %v", err)
	}
	defer minifluxSecretFileHandle.Close()
	configFileBytes, err := ioutil.ReadAll(minifluxSecretFileHandle)
	if err != nil {
		return nil, fmt.Errorf("Miniflux Config file could not be read: %v", err)
	}
	var config minifluxConfig
	err = json.Unmarshal(configFileBytes, &config)
	if err != nil || config.Url == "" || config.Token == "" {
		return nil, fmt.Errorf("Miniflux Config file could not been understood")
	}
	
	// create client
	client := &MinifluxClient{miniflux.New(config.Url, config.Token), nil}

	// find youtube Category
	categories, err := client.client.Categories()
	if err != nil {
		return nil, err
	}
	for _, category := range categories {
		if category.Title == "YouTube" {
			client.youtubeCategory = category
			break
		}
	} 
	if client.youtubeCategory == nil {
		category, err := client.client.CreateCategory("YouTube")
		if err != nil {
			return nil, fmt.Errorf("Category YouTube does not exist and creation failed: %v", err)
		}
		client.youtubeCategory = category
	}

	return client, nil
}

func (client *MinifluxClient) GetYoutubeSubscriptions() ([]miniflux.Feed, error) {
	feeds, err := client.client.Feeds()
	if err != nil {
		return nil, err
	}

	results := []miniflux.Feed{}

	for _, feed := range feeds {
		if feed.Category.ID == client.youtubeCategory.ID {
			results = append(results, *feed)
		}
	}

	return results, nil
}

func (client *MinifluxClient) Subscribe(feedURL string) error {
	_, err := client.client.CreateFeed(feedURL, client.youtubeCategory.ID)
	return err
}

func (client *MinifluxClient) Unsubscribe(feed miniflux.Feed) error {
	return client.client.DeleteFeed(feed.ID)
}
//Lib functions for subscribe and unsubscribe
//func (c *Client) CreateFeed(url string, categoryID int64) (int64, error)
//func (c *Client) DeleteFeed(feedID int64) error