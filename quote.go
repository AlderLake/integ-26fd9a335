/*
Package quote is free quote downloader library and cli

Downloads daily/weekly/monthly historical price quotes from Yahoo
and daily/intraday data from Tiingo/Bittrex/Binance

Copyright 2019 Mark Chenoweth
Licensed under terms of MIT license (see LICENSE)
*/
package quote

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/textproto"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Quote - stucture for historical price data
type Quote struct {
	Symbol    string      `json:"symbol"`
	Precision int64       `json:"-"`
	Date      []time.Time `json:"date"`
	Open      []float64   `json:"open"`
	High      []float64   `json:"high"`
	Low       []float64   `json:"low"`
	Close     []float64   `json:"close"`
	Volume    []float64   `json:"volume"`
}

// Quotes - an array of historical price data
type Quotes []Quote

// Period - for quote history
type Period string

// ClientTimeout - connect/read timeout for client requests
const ClientTimeout = 10 * time.Second

const (
	// Min1 - 1 Minute time period
	Min1 Period = "60"
	// Min3 - 3 Minute time period
	Min3 Period = "3m"
	// Min5 - 5 Minute time period
	Min5 Period = "300"
	// Min15 - 15 Minute time period
	Min15 Period = "900"
	// Min30 - 30 Minute time period
	Min30 Period = "1800"
	// Min60 - 60 Minute time period
	Min60 Period = "3600"
	// Hour2 - 2 hour time period
	Hour2 Period = "2h"
	// Hour4 - 4 hour time period
	Hour4 Period = "4h"
	// Hour6 - 6 hour time period
	Hour6 Period = "6h"
	// Hour8 - 8 hour time period
	Hour8 Period = "8h"
	// Hour12 - 12 hour time period
	Hour12 Period = "12h"
	// Daily time period
	Daily Period = "d"
	// Day3 - 3 day time period
	Day3 Period = "3d"
	// Weekly time period
	Weekly Period = "w"
	// Monthly time period
	Monthly Period = "m"
)

// Log - standard logger, disabled by default
var Log *log.Logger

// Delay - time delay in milliseconds between quote requests (default=100)
// Be nice, don't get blocked
var Delay time.Duration

func init() {
	Log = log.New(ioutil.Discard, "quote: ", log.Ldate|log.Ltime|log.Lshortfile)
	Delay = 100
}

// NewQuote - new empty Quote struct
func NewQuote(symbol string, bars int) Quote {
	return Quote{
		Symbol: symbol,
		Date:   make([]time.Time, bars),
		Open:   make([]float64, bars),
		High:   make([]float64, bars),
		Low:    make([]float64, bars),
		Close:  make([]float64, bars),
		Volume: make([]float64, bars),
	}
}

// ParseDateString - parse a potentially partial date string to Time
func ParseDateString(dt string) time.Time {
	if dt == "" {
		return time.Now()
	}
	t, _ := time.Parse("2006-01-02 15:04", dt+"0000-01-01 00:00"[len(dt):])
	return t
}

func getPrecision(symbol string) int {
	var precision int
	precision = 2
	if strings.Contains(strings.ToUpper(symbol), "BTC") ||
		strings.Contains(strings.ToUpper(symbol), "ETH") ||
		strings.Contains(strings.ToUpper(symbol), "USD") {
		precision = 8
	}
	return precision
}

// CSV - convert Quote structure to csv string
func (q Quote) CSV() string {

	precision := getPrecision(q.Symbol)

	var buffer bytes.Buffer
	buffer.WriteString("datetime,open,high,low,close,volume\n")
	for bar := range q.Close {
		str := fmt.Sprintf("%s,%.*f,%.*f,%.*f,%.*f,%.*f\n", q.Date[bar].Format("2006-01-02 15:04"),
			precision, q.Open[bar], precision, q.High[bar], precision, q.Low[bar], precision, q.Close[bar], precision, q.Volume[bar])
		buffer.WriteString(str)
	}
	return buffer.String()
}

// Highstock - convert Quote structure to Highstock json format
func (q Quote) Highstock() string {

	precision := getPrecision(q.Symbol)

	var buffer bytes.Buffer
	buffer.WriteString("[\n")
	for bar := range q.Close {
		comma := ","
		if bar == len(q.Close)-1 {
			comma = ""
		}
		str := fmt.Sprintf("[%d,%.*f,%.*f,%.*f,%.*f,%.*f]%s\n",
			q.Date[bar].UnixNano()/1000000, precision, q.Open[bar], precision, q.High[bar], precision, q.Low[bar], precision, q.Close[bar], precision, q.Volume[bar], comma)
		buffer.WriteString(str)

	}
	buffer.WriteString("]\n")
	return buffer.String()
}

// Amibroker - convert Quote structure to csv string
func (q Quote) Amibroker() string {

	precision := ge