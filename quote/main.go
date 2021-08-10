
/*
Package quote is free quote downloader library and cli

Downloads daily/weekly/monthly/yearly historical price quotes from Yahoo
and daily/intraday data from Tiingo, crypto from Coinbase/Bittrex/Binance

Copyright 2019 Mark Chenoweth
Licensed under terms of MIT license

*/

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/markcheno/go-quote"
)

var usage = `Usage:
  quote -h | -help
  quote -v | -version
  quote <market> [-output=<outputFile>]
  quote [-years=<years>|(-start=<datestr> [-end=<datestr>])] [options] [-infile=<filename>|<symbol> ...]

Options:
  -h -help             show help
  -v -version          show version
  -years=<years>       number of years to download [default=5]
  -start=<datestr>     yyyy[-[mm-[dd]]]
  -end=<datestr>       yyyy[-[mm-[dd]]] [default=today]
  -infile=<filename>   list of symbols to download
  -outfile=<filename>  output filename
  -period=<period>     1m|3m|5m|15m|30m|1h|2h|4h|6h|8h|12h|d|3d|w|m [default=d]
  -source=<source>     yahoo|tiingo|tiingo-crypto|coinbase|bittrex|binance [default=yahoo]
  -token=<tiingo_tok>  tingo api token [default=TIINGO_API_TOKEN]
  -format=<format>     (csv|json|hs|ami) [default=csv]
  -adjust=<bool>       adjust yahoo prices [default=true]
  -all=<bool>          all in one file (true|false) [default=false]
  -log=<dest>          filename|stdout|stderr|discard [default=stdout]
  -delay=<ms>          delay in milliseconds between quote requests

Note: not all periods work with all sources

Valid markets:
etfs:       etf
crypto:     bittrex-btc,bittrex-eth,bittrex-usdt,
            binance-bnb,binance-btc,binance-eth,binance-usdt,
            coinbase
`

const (
	version    = "0.2"
	dateFormat = "2006-01-02"
)

type quoteflags struct {
	years   int
	delay   int
	start   string
	end     string
	period  string
	source  string
	token   string
	infile  string
	outfile string
	format  string
	log     string
	all     bool
	adjust  bool
	version bool
}

func check(e error) {
	if e != nil {
		fmt.Printf("\nerror: %v\n\n", e)
		fmt.Println(usage)
		os.Exit(0)
		//panic(e)
	}
}

func checkFlags(flags quoteflags) error {

	// validate source
	if flags.source != "yahoo" &&
		flags.source != "tiingo" &&
		flags.source != "tiingo-crypto" &&
		flags.source != "coinbase" &&
		flags.source != "bittrex" &&
		flags.source != "binance" {
		return fmt.Errorf("invalid source, must be either 'yahoo', 'tiingo', 'coinbase', 'bittrex', or 'binance'")
	}

	// validate period
	if flags.source == "yahoo" &&
		(flags.period == "1m" || flags.period == "5m" || flags.period == "15m" || flags.period == "30m" || flags.period == "1h") {
		return fmt.Errorf("invalid period for yahoo, must be 'd'")
	}
	if flags.source == "tiingo" {
		// check period
		if flags.period != "d" {
			return fmt.Errorf("invalid period for tiingo, must be 'd'")
		}
		// check token
		if flags.token == "" {
			return fmt.Errorf("missing token for tiingo, must be passed or TIINGO_API_TOKEN must be set")
		}
	}

	if flags.source == "tiingo-crypto" &&
		!(flags.period == "1m" ||
			flags.period == "3m" ||
			flags.period == "5m" ||
			flags.period == "15m" ||
			flags.period == "30m" ||
			flags.period == "1h" ||
			flags.period == "2h" ||
			flags.period == "4h" ||
			flags.period == "6h" ||
			flags.period == "8h" ||
			flags.period == "12h" ||
			flags.period == "d") {
		return fmt.Errorf("invalid source for tiingo-crypto, must be '1m', '3m', '5m', '15m', '30m', '1h', '2h', '4h', '6h', '8h', '12h', '1d', '3d', '1w', or '1M'")
	}

	if flags.source == "tiingo-crypto" && flags.token == "" {
		return fmt.Errorf("missing token for tiingo-crypto, must be passed or TIINGO_API_TOKEN must be set")
	}

	if flags.source == "bittrex" && !(flags.period == "1m" || flags.period == "5m" || flags.period == "30m" || flags.period == "1h" || flags.period == "d") {
		return fmt.Errorf("invalid source for bittrex, must be '1m', '5m', '30m', '1h' or 'd'")
	}

	if flags.source == "binance" &&
		!(flags.period == "1m" ||
			flags.period == "3m" ||
			flags.period == "5m" ||
			flags.period == "15m" ||
			flags.period == "30m" ||
			flags.period == "1h" ||
			flags.period == "2h" ||
			flags.period == "4h" ||
			flags.period == "6h" ||
			flags.period == "8h" ||
			flags.period == "12h" ||
			flags.period == "d" ||
			flags.period == "3d" ||
			flags.period == "w" ||
			flags.period == "m") {