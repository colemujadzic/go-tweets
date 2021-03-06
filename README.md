# go-tweets

[![Travis CI](https://img.shields.io/travis/colemujadzic/go-tweets.svg?style=for-the-badge)](https://travis-ci.org/colemujadzic/go-tweets)

View a twitter account's most recent tweets with this command-line tool, written in Go

## Installation

#### Via Go

```console
$ go get github.com/colemujadzic/go-tweets
```

## Usage

```console
$ go-tweets -h
go-tweets - A command line tool to retrieve a twitter user's tweets using Twitter's API

Version: 0.0.1

Usage:  go-tweets [options] <twitter username> <number of tweets>

  -consumer-key string
        twitter consumer key
  -consumer-secret string
        twitter consumer secret
```

## Acknowledgments

* Most of the inspiration to try this came from [Jessie Frazelle's](https://github.com/jessfraz/) command-line tools.
* [Filippo Valsorda's](https://github.com/FiloSottile) amazing [Makefile](https://github.com/cloudflare/hellogopher) for Go was also super helpful.
