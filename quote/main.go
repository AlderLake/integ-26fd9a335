
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
		return fmt.Errorf("invalid source for binance, must be '1m', '3m', '5m', '15m', '30m', '1h', '2h', '4h', '6h', '8h', '12h', '1d', '3d', '1w', or '1M'")
	}

	return nil
}

func setOutput(flags quoteflags) error {
	var err error
	if flags.log == "stdout" {
		quote.Log.SetOutput(os.Stdout)
	} else if flags.log == "stderr" {
		quote.Log.SetOutput(os.Stderr)
	} else if flags.log == "discard" {
		quote.Log.SetOutput(ioutil.Discard)
	} else {
		var f *os.File
		f, err = os.OpenFile(flags.log, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		defer f.Close()
		quote.Log.SetOutput(f)
	}
	return err
}

func getSymbols(flags quoteflags, args []string) ([]string, error) {

	var err error
	var symbols []string

	if flags.infile != "" {
		symbols, err = quote.NewSymbolsFromFile(flags.infile)
		if err != nil {
			return symbols, err
		}
	} else {
		symbols = args
	}

	// make sure we found some symbols
	if len(symbols) == 0 {
		return symbols, fmt.Errorf("no symbols specified")
	}

	// validate outfileFlag
	if len(symbols) > 1 && flags.outfile != "" && !flags.all {
		return symbols, fmt.Errorf("outfile not valid with multiple symbols\nuse -all=true")
	}

	return symbols, nil
}

func getPeriod(periodFlag string) quote.Period {
	period := quote.Daily
	switch periodFlag {
	case "1m":
		period = quote.Min1
	case "3m":
		period = quote.Min3
	case "5m":
		period = quote.Min5
	case "15m":
		period = quote.Min15
	case "30m":
		period = quote.Min30
	case "1h":
		period = quote.Min60
	case "2h":
		period = quote.Hour2
	case "4h":
		period = quote.Hour4
	case "6h":
		period = quote.Hour6
	case "8h":
		period = quote.Hour8
	case "12h":
		period = quote.Hour12
	case "d":
		period = quote.Daily
	case "1d":
		period = quote.Daily
	case "3d":
		period = quote.Day3
	case "w":
		period = quote.Weekly
	case "1w":
		period = quote.Weekly
	case "m":
		period = quote.Monthly
	case "1M":
		period = quote.Monthly
	}
	return period
}

func getTimes(flags quoteflags) (time.Time, time.Time) {
	// determine start/end times
	to := quote.ParseDateString(flags.end)
	var from time.Time
	if flags.start != "" {
		from = quote.ParseDateString(flags.start)
	} else { // use years
		from = to.Add(-time.Duration(int(time.Hour) * 24 * 365 * flags.years))
	}
	return from, to
}

func outputAll(symbols []string, flags quoteflags) error {
	// output all in one file
	from, to := getTimes(flags)
	period := getPeriod(flags.period)
	quotes := quote.Quotes{}
	var err error
	if flags.source == "yahoo" {
		quotes, err = quote.NewQuotesFromYahooSyms(symbols, from.Format(dateFormat), to.Format(dateFormat), period, flags.adjust)
	} else if flags.source == "tiingo" {
		quotes, err = quote.NewQuotesFromTiingoSyms(symbols, from.Format(dateFormat), to.Format(dateFormat), flags.token)
	} else if flags.source == "tiingo-crypto" {
		quotes, err = quote.NewQuotesFromTiingoCryptoSyms(symbols, from.Format(dateFormat), to.Format(dateFormat), period, flags.token)
	} else if flags.source == "coinbase" {
		quotes, err = quote.NewQuotesFromCoinbaseSyms(symbols, from.Format(dateFormat), to.Format(dateFormat), period)
	} else if flags.source == "bittrex" {
		quotes, err = quote.NewQuotesFromBittrexSyms(symbols, period)
	} else if flags.source == "binance" {
		quotes, err = quote.NewQuotesFromBinanceSyms(symbols, from.Format(dateFormat), to.Format(dateFormat), period)
	}
	if err != nil {
		return err
	}

	if flags.format == "csv" {
		err = quotes.WriteCSV(flags.outfile)
	} else if flags.format == "json" {
		err = quotes.WriteJSON(flags.outfile, false)
	} else if flags.format == "hs" {
		err = quotes.WriteHighstock(flags.outfile)
	} else if flags.format == "ami" {
		err = quotes.WriteAmibroker(flags.outfile)
	}
	return err