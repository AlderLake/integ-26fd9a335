# go-quote

[![GoDoc](http://godoc.org/github.com/markcheno/go-quote?status.svg)](http://godoc.org/github.com/markcheno/go-quote) 

A free quote downloader library and cli 

Downloads daily historical price quotes from Yahoo and daily/intraday data from various api's. Written in pure Go. No external dependencies. Now downloads crypto coin historical data from various exchanges.

- Update: 11/15/2021 - Removed obsolete markets, converted to go modules

- Update: 7/18/2021 - Removed obsolete Google support

- Update: 6/26/2019 - updated GDAX to Coinbase, added coinbase market

- Update: 4/26/2018 - Added preliminary [tiingo](https://api.tiingo.com/) CRYPTO support. Use -source=tiingo-crypto -token=<your_tingo_token> You can also set env variable TIINGO_API_TOKEN. To get symbol lists, use market: tiingo-btc, tiingo-eth or tiingo-usd

- Update: 12/21/2017 - Added Amibroker format option (creates csv file with separate date and time). Use -format=ami

- Update: 12/20/2017 - Added [Binance](https://www.binance.com/trade.html) exchange support. Use -source=binance

- Update: 12/18/2017 - Added [Bittrex](https://bittrex.com/home/markets) exchange support. Use -source=bittrex  

- Update: 10/21/2017 - Added Coinbase [GDAX](https://www.gdax.com/trade/BTC-USD) exchange support. Use -source=gdax All times are in UTC. Automatically rate limited. 

- Update: 7/19/2017 - Added preliminary [tiingo](https://api.tiingo.com/) support. Use -source=tiingo -token=<your_tingo_token> You can also set env variable TIINGO_API_TOKEN

- Update: 5/24/2017 - Now works with the new Yahoo download format. Beware - Yahoo data quality is now questionable and the free Yahoo quotes are likely to permanently go away in the near future. Use with caution!

Still very much in alpha mode. Expect bugs and API changes. Comments/suggestions/pull requests welcome!

Copyright 2018 Mark Chenoweth

Install CLI utility (quote) with:

```bash
go install github.com/markcheno/go-quote/quote@latest
```

```
Usage:
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
  -infile=<filename>   lis