package main

import (
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"net/http"
	"os"
	"runtime"
	"time"
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
	r.Header.Set("User-Agent", runtime.GOOS+":io.anglur.GRNotifier:1.0 (by /u/"+config.Username+")")
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

func parseRSSFeed(config Config, feed *gofeed.Feed) {
	reset := false

	for _, item := range feed.Items {
		t, e := time.Parse(time.RFC3339, item.Updated)

		if e != nil {
			fmt.Println(e)
			continue
		}

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

	if reset {
		start = time.Now().Unix()
	}
}

var config = LoadConfig()
var start = time.Now().Unix()

func main() {
	fp := gofeed.NewParser()
	fp.Client = &http.Client{
		Transport: &UserAgentTransport{http.DefaultTransport},
	}

	for {
		for _, subreddit := range config.SubReddits {
			feed, e := fp.ParseURL(config.BaseURL + "/r/" + subreddit.URL + "/new/.rss")

			if e != nil {
				fmt.Println(e)
				continue
			}

			parseRSSFeed(config, feed)
		}

		fmt.Printf("Sleeping for %d seconds...\n", config.Interval)
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
}

func sendNotification(n Notification) {
	fmt.Println(n.Title + " - " + n.Time + " \n" + n.URL + " \n" + n.Author + " \n" + n.Subreddit + "\n")
}
