
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