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
	BaseURL    string         `json:"baseurl"`
	Interval   int         `json:"interval"`
	Username   string      `json:"username"`
	SubReddits []SubReddit `json:"subreddits"`
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

func parseRSSFeed(config Config, start int64, feed *gofeed.Feed) {
	reset := false
	for _, item := range feed.Items {
		t, _ := time.Parse(time.RFC3339, item.Updated)
		if start < t.Unix() {
			title := item.Title
			link := item.Link
			author := item.Author.Name
			subreddit := item.Categories[0]

			fmt.Println(title + " - " + t.String() + " \n" + link + " \n" + config.BaseURL + author + " \n" + config.BaseURL + "/r/" + subreddit + "\n")

			reset = true
		}
	}

	if reset {
		start = time.Now().Unix()
	}
}

func main() {
	config := LoadConfig()

	start := time.Now().Unix()

	fp := gofeed.NewParser()
	fp.Client = &http.Client{
		Transport: &UserAgentTransport{http.DefaultTransport},
	}

	fmt.Println(config.SubReddits)


	for _, subreddit := range config.SubReddits {
		feed, e := fp.ParseURL(config.BaseURL + "/r/" + subreddit.URL + "/new/.rss")

		if e != nil {
			fmt.Println(e)
		}

		parseRSSFeed(config, start, feed)
	}
}
