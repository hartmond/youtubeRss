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

// YoutubeClient wraps the YouTube API client for simplified use
type YoutubeClient struct {
	client *youtube.Service
}

// NewYoutubeClient initializes a YouTube API client with configurations from json files
func NewYoutubeClient(youtubeSecretFile, youtubeUserSecretFile string) (*YoutubeClient, error) {
	scope := youtube.YoutubeReadonlyScope
	ctx := context.Background()

	b, err := ioutil.ReadFile(youtubeSecretFile)
	if err != nil {
		return nil, fmt.Errorf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, scope)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse client secret file to config: %v", err)
	}

	var token *oauth2.Token
	if token, err = loadOauthToken(youtubeUserSecretFile); err != nil {
		token, err = setupOauth(config, youtubeUserSecretFile)
		if err != nil {
			return nil, fmt.Errorf("Oauth setup failed: %v", err)
		}
	}

	httpClient := config.Client(ctx, token)
	youtubeClient, err := youtube.New(httpClient)
	if err != nil {
		return nil, err
	}

	return &YoutubeClient{youtubeClient}, nil
}

func loadOauthToken(youtubeUserSecretFile string) (*oauth2.Token, error) {
	f, err := os.Open(youtubeUserSecretFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

func setupOauth(config *oauth2.Config, youtubeUserSecretFile string) (*oauth2.Token, error) {
	fmt.Println("Starting Oauth setup as user configuration is not available or invalid")

	config.RedirectURL = "urn:ietf:wg:oauth:2.0:oob"
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	var code string
	fmt.Printf("Open the following link in your browser:\n\n%v\n\n", authURL)
	fmt.Println("Enter the displayed token:")

	if _, err := fmt.Scan(&code); err != nil {
		return nil, fmt.Errorf("Unable to read authorization code %v", err)
	}

	fmt.Println(authURL)

	token, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve token %v", err)
	}

	fmt.Println("Saving token")
	fmt.Printf("Saving credential file to: %s\n", youtubeUserSecretFile)
	f, err := os.OpenFile(youtubeUserSecretFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, fmt.Errorf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)

	return token, nil
}

// GetSubscriptions returns a list of the FeedURLs of all currently subscribed YouTube channels
func (client *YoutubeClient) GetSubscriptions() ([]string, error) {
	results := []string{}

	for pageToken := ""; ; {
		res, err := client.client.Subscriptions.List("snippet").Mine(true).MaxResults(10).PageToken(pageToken).Do()
		if err != nil {
			return nil, err
		}

		for _, elem := range res.Items {
			results = append(results, fmt.Sprintf("https://www.youtube.com/feeds/videos.xml?channel_id=%s", elem.Snippet.ResourceId.ChannelId))
		}

		pageToken = res.NextPageToken
		if pageToken == "" {
			return results, nil
		}
	}
}
