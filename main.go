package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/kurrik/twittergo"

	"github.com/kurrik/oauth1a"
)

const (
	// banner to be displayed
	BANNER = `
Go-Tweets

Get tweets for a user
Version: %s

`
	// version
	VERSION = "0.0.1"
)

var (
	version               bool
	twitterUser           string
	twitterConsumerKey    string
	twitterConsumerSecret string
)

func init() {
	flag.StringVar(&twitterConsumerKey, "consumer-key", "", "twitter consumer key")
	flag.StringVar(&twitterConsumerSecret, "consumer-secret", "", "twitter consumer secret")

	flag.BoolVar(&version, "version", false, "print version and exit")
	flag.BoolVar(&version, "v", false, "print version and exit (s)")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, VERSION))
		flag.PrintDefaults()
	}

	flag.Parse()

	if twitterConsumerKey == "" {
		if twitterConsumerKey == os.Getenv("CONSUMER_KEY") {
			twitterConsumerKey = ""
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	if twitterConsumerSecret == "" {
		if twitterConsumerSecret == os.Getenv("CONSUMER_SECRET") {
			twitterConsumerSecret = ""
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	if len(os.Args) < 1 {
		fmt.Println("Provide a username, e.g. @dril")
		os.Exit(1)
	}

	argument := flag.Args()[0]

	if argument == "help" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	twitterUser = argument

}

func main() {
	// create config
	config := &oauth1a.ClientConfig{
		ConsumerKey:    "CONSUMER_KEY",
		ConsumerSecret: "CONSUMER_SECRET",
	}

	client := twittergo.NewClient(config, nil)
	if err := client.FetchAppToken(); err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't fetch app token: %v\n", err)
		os.Exit(2)
	}

	// don't save app token
	_ = client.GetAppToken()

	// make request
	value := url.Values{}
	value.Set("user_id", twitterUser)
	request, err := http.NewRequest("GET", "/1.1/statuses/user_timeline.json"+value.Encode(), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't parse request: %v\n", err)
		os.Exit(2)
	}

	response, err := client.SendRequest(request)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't send request: %v\n", err)
		os.Exit(2)
	}

	// get response
	searchResults := &twittergo.SearchResults{}
	if err := response.Parse(searchResults); err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't parse response: %v\n", err)
		os.Exit(2)
	}

	// print tweets
	tweets := searchResults.Statuses()

	for _, value := range tweets {
		fmt.Println(value.Text())
	}

}
