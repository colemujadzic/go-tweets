package main

import (
	"flag"
	"fmt"
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

	if len(os.Args) < 2 {
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
	config := &oauth1a.ClientConfig{
		ConsumerKey:    "CONSUMER_KEY",
		ConsumerSecret: "CONSUMER_SECRET",
	}

	client := twittergo.NewClient(config, nil)
	if err := client.FetchAppToken(); err != nil {
		// handle err
	}
}
