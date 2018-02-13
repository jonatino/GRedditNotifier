package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"time"
	"net/http"
	"os"
	"encoding/json"
)

type SubReddit struct {
	URL         string `json:"url"`
	IFTTTApiKey string `json:"ifttt_api_key"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phonenumber"`
}

type Config struct {
	BaseURL    string      `json:"baseurl"`
	Interval   int         `json:"interval"`
	Username   string      `json:"username"`
	SubReddits []SubReddit `json:"subreddits"`
}

type Notification struct {
	Title     string
	Time      string
	URL       string
	Author    string
	Subreddit string
}

type UserAgentTransport struct {
	http.RoundTripper
}

func (c *UserAgentTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("User-Agent", "windows:io.anglur.GRNotifier:1.0 (by /u/Tiflotin)")
	return c.RoundTripper.RoundTrip(r)
}

func LoadConfig() Config {
	var config Config
	configFile, err := os.Open("./GNotifierConfig.json")
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}

func parseRSSFeed(config Config, start int64, feed *gofeed.Feed) bool {
	//TODO possible to pass start by pointer so i don't have to return value?
	reset := false

	for _, item := range feed.Items {
		t, _ := time.Parse(time.RFC3339, item.Updated)
		item.Updated = t.String()

		if start < t.Unix() {
			notification := Notification{
				Title:     item.Title,
				Time:      item.Updated,
				URL:       item.Link,
				Author:    config.BaseURL + item.Author.Name,
				Subreddit: config.BaseURL + "/r/" + item.Categories[0],
			}

			sendNotification(notification)

			reset = true
		}
	}

	return reset
}

func main() {
	config := LoadConfig()
	start := time.Now().Unix() - 200000

	fp := gofeed.NewParser()
	fp.Client = &http.Client{
		Transport: &UserAgentTransport{http.DefaultTransport},
	}

	reset := false

	for _, subreddit := range config.SubReddits {
		feed, e := fp.ParseURL(config.BaseURL + "/r/" + subreddit.URL + "/new/.rss")

		if e != nil {
			fmt.Println(e)
		}

		reset = parseRSSFeed(config, start, feed)
	}

	if reset {
		start = time.Now().Unix()
	}

}

func sendNotification(n Notification) {
	fmt.Println(n.Title + " - " + n.Time + " \n" + n.URL + " \n" + n.Author + " \n" + n.Subreddit + "\n")
}
