package main

import (
	"fmt"
	miniflux "miniflux.app/client"
)

type MinifluxClient struct {
	client          *miniflux.Client
	youtubeCategory *miniflux.Category
}

func NewMinifluxClient(url, token string) (*MinifluxClient, error) {
	client := &MinifluxClient{}
	client.client = miniflux.New(url, token)

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
		return nil, fmt.Errorf("Category YouTube not found")
		// TODO: create category
	}

	return client, nil
}

func (client *MinifluxClient) GetYoutubeSubscriptions() ([]string, error) {
	feeds, err := client.client.Feeds()
	if err != nil {
		return nil, err
	}

	results := []string{}

	for _, feed := range feeds {
		if feed.Category.ID == client.youtubeCategory.ID {
			results = append(results, feed.FeedURL)
		}
	}

	return results, nil
}

//Lib functions for subscribe and unsubscribe
//func (c *Client) CreateFeed(url string, categoryID int64) (int64, error)
//func (c *Client) DeleteFeed(feedID int64) error