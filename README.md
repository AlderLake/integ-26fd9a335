# go-quote

[![GoDoc](http://godoc.org/github.com/markcheno/go-quote?status.svg)](http://godoc.org/github.com/markcheno/go-quote) 

A free quote downloader library and cli 

Downloads daily historical price quotes from Yahoo and daily/intraday data from various api's. Written in pure Go. No external dependencies. Now downloads crypto coin historical data from various exchanges.

- Update: 11/15/2021 - Removed obsolete markets, converted to go modules

- Update: 7/18/2021 - Removed obsolete Google support

- Update: 6/26/2019 - updated GDAX to Coinbase, added coinbase market

- Update: 4/26/2018 - Added preliminary [tiingo](https://api.tiingo.com/) CRYPTO support. Use -source=tiingo-crypto -token=<your_ting