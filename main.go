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

Version: %s
`
	// VERSION ...
	VERSION = "0.0.1"
)

var (
	version               bool
	tweet                 []byte
	twitterUser           string
	numberOfTweets        string
	twitterConsumerKey    string
	twitterConsumerSecret string
)

func init() {
	flag.StringVar(&twitterConsumerKey, "consumer-key", "", "Twitter consumer key")
	flag.StringVar(&twitterConsumerSecret, "consumer-secret", "", "Twitter consumer secret")

	flag.BoolVar(&version, "version", false, "Print version and exit")
	flag.BoolVar(&version, "v", false, "Print version and exit")

	// flag.Usage prints usage information when -h or -help flag is invoked
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(BANNER, VERSION))
		fmt.Println()
		fmt.Println("Usage:  go-tweets  [options]  <twitter username>  <number of tweets>")
		fmt.Println()
		flag.PrintDefaults()
	}

	flag.Parse()

	if twitterConsumerKey == "" {
		if twitterConsumerKey = os.Getenv("CONSUMER_KEY"); twitterConsumerKey == "" {
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	if twitterConsumerSecret == "" {
		if twitterConsumerSecret = os.Getenv("CONSUMER_SECRET"); twitterConsumerSecret == "" {
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	/*
		// reserve for commands / subcommands
		// verify that a subommand has been provided
		// os.Args[0] is the main command
		// os.Args[2] is the subcommand
		if len(os.Args) == 0 {
			fmt.Println("A command or subcommand is required")
			flag.PrintDefaults()
			os.Exit(1)
		}
	*/

	// check for arguments
	if flag.NArg() == 0 {
		fmt.Println("Please provide a username.")
		flag.PrintDefaults()
		os.Exit(1)
	} else if flag.NArg() == 1 {
		fmt.Println("Please provide the number of tweets you wish to retrieve.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// argument := flag.Args()[0]
	firstArgument := flag.Args()[0]
	secondArgument := flag.Args()[1]

	// note:

	if firstArgument == "version" || firstArgument == "v" {
		fmt.Println("VERSION - TEST")
		fmt.Printf("%s", VERSION)
		os.Exit(0)
	}

	twitterUser = firstArgument
	numberOfTweets = secondArgument

	if version {
		fmt.Printf("%s", VERSION)
		os.Exit(0)
	}
}

func main() {
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
