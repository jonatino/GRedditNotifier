package main

import (
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/xconstruct/go-pushbullet"
	"net/http"
	"os"
	"runtime"
	"time"
)

type SubReddit struct {
	URL string `json:"url"`
}

type Config struct {
	BaseURL          string      `json:"baseurl"`
	Interval         int         `json:"interval"`
	Username         string      `json:"username"`
	PushBulletApiKey string      `json:"pushbullet_api_key"`
	SubReddits       []SubReddit `json:"subreddits"`
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
				Subreddit: "/r/" + item.Categories[0],
			}

			sendNotification(notification)

			reset = true
		}
	}

	if reset {
		start = time.Now().Unix()
	}
}

func LoadPushBulletDevice() *pushbullet.Device {
	devices, e := pb.Devices()
	if e != nil {
		panic(e)
	}

	var device *pushbullet.Device
	for _, dev := range devices {
		if dev.Active == true {
			device = dev
			break
		}
	}

	if device == nil {
		fmt.Println("Could not find active device on PushBullet!")
		os.Exit(2)
	}

	return device
}

var config = LoadConfig()
var start = time.Now().Unix()
var pb = pushbullet.New(config.PushBulletApiKey)
var device = LoadPushBulletDevice()

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
	body := n.Title + " - " + n.Time + " \n" + n.URL + " \n" + n.Author + " \n" + n.Subreddit + "\n"

	fmt.Println(body)

	e := pb.PushNote(device.Iden, "New "+n.Subreddit+" post!", body)

	if e != nil {
		panic(e)
	}

}
