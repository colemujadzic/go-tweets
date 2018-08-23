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
	tweet                 []byte
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
		if twitterConsumerKey = os.Getenv("CONSUMER_KEY"); twitterConsumerKey == "" {
			fmt.Println("TEST: COULDN'T GET CONSUMER KEY")
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	if twitterConsumerSecret == "" {
		if twitterConsumerSecret = os.Getenv("CONSUMER_SECRET"); twitterConsumerSecret == "" {
			fmt.Println("TEST: COULDN'T GET CONSUMER SECRET")
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	if flag.NArg() == 0 {
		fmt.Println("TEST: NO USER ENTERED")
		flag.PrintDefaults()
		os.Exit(1)
	}

	argument := flag.Args()[0]

	if argument == "help" {
		fmt.Println("TEST: HELP ARGUMENT")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if argument == "version" {
		fmt.Println("TEST: VERSION ARGUMENT")
		flag.PrintDefaults()
		os.Exit(1)
	}

	twitterUser = argument

	if version {
		fmt.Printf("%s", VERSION)
		os.Exit(0)
	}
}

func main() {
	numberOfTweets := "1"
	// create config
	config := &oauth1a.ClientConfig{
		ConsumerKey:    twitterConsumerKey,
		ConsumerSecret: twitterConsumerSecret,
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
	value.Set("count", numberOfTweets)
	value.Set("screen_name", twitterUser)
	request, err := http.NewRequest("GET", "/1.1/statuses/user_timeline.json?"+value.Encode(), nil)
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
