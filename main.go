package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
)

const (
	// BANNER ...
	BANNER = `
go-tweets - Command line tool to retrieve a twitter user's tweets using Twitter's API
`
)

var (
	tweet                 []byte
	twitterUser           string
	numberOfTweets        string
	twitterConsumerKey    string
	twitterConsumerSecret string
)

func init() {
	// flag definitions
	flag.StringVar(&twitterConsumerKey, "consumer-key", os.Getenv("TWITTER_CONSUMER_KEY"), "Twitter consumer key (or env var TWITTER_CONSUMER_KEY)")
	flag.StringVar(&twitterConsumerSecret, "consumer-secret", os.Getenv("TWITTER_CONSUMER_SECRET"), "Twitter consumer secret (or env var TWITTER_CONSUMER_SECRET)")

	// flag.Usage prints usage information when -h or -help flag is invoked
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER))
		fmt.Println()
		fmt.Println("Usage:  go-tweets  [options]  <twitter username>  <number of tweets>")
		fmt.Println()
		flag.PrintDefaults()
	}

	// parse command-line flags
	flag.Parse()

	// check for consumer-key flag
	if twitterConsumerKey == "" {
		if twitterConsumerKey = os.Getenv("TWITTER_CONSUMER_KEY"); twitterConsumerKey == "" {
			fmt.Println("Please provide a Twitter API consumer key")
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	// check for consumer-secret flag
	if twitterConsumerSecret == "" {
		if twitterConsumerSecret = os.Getenv("TWITTER_CONSUMER_SECRET"); twitterConsumerSecret == "" {
			fmt.Println("Please provide a Twitter API consumer secret")
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	// assign command-line arguments
	firstArgument := flag.Args()[0]
	secondArgument := flag.Args()[1]

	// check for arguments that remain after parsing flags
	if flag.NArg() == 0 {
		fmt.Println("Please provide a username.")
		flag.PrintDefaults()
		os.Exit(1)
	} else if flag.NArg() == 1 {
		fmt.Println("Please provide the number of tweets you wish to retrieve.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// assign arguments to variables
	twitterUser = firstArgument
	numberOfTweets = secondArgument
}

func main() {
	// create config
	config := &oauth1a.ClientConfig{
		ConsumerKey:    twitterConsumerKey,
		ConsumerSecret: twitterConsumerSecret,
	}

	// create new twitter client
	client := twittergo.NewClient(config, nil)
	if err := client.FetchAppToken(); err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't fetch app token: %v\n", err)
		os.Exit(2)
	}

	// let's not save the app token
	_ = client.GetAppToken()

	// make new request
	value := url.Values{}
	value.Set("count", numberOfTweets)
	value.Set("screen_name", twitterUser)
	request, err := http.NewRequest("GET", "/1.1/statuses/user_timeline.json?"+value.Encode(), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't parse request: %v\n", err)
		os.Exit(2)
	}

	// send request
	response, err := client.SendRequest(request)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't send request: %v\n", err)
		os.Exit(2)
	}

	// parse response
	results := &twittergo.Timeline{}
	if err := response.Parse(results); err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't parse response: %v\n", err)
		os.Exit(2)
	}

	// print tweets
	for _, value := range *results {
		if tweet, err = json.Marshal(*results); err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't encode tweet: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(value.Text())
	}
}
