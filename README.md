# Yahoo Finance CLI Tool

## Description 

A command-line utility tool for downloading historical stock data from Yahoo Finance. 

### [Yahoo Finance API](https://cryptocointracker.com/yahoo-finance/yahoo-finance-api)

* [https://query2.finance.yahoo.com/v8/finance/chart/{symbol}]()
  * download stock ticker
  * interval: "1d","5d","1mo","3mo","6mo","1y","2y","5y","10y","max"
  * period1: start time - UNIX timestamp
  * period2: end time - UNIX timestamp

#### Yahoo Finance API Request Parameters

```
https://query2.finance.yahoo.com/v8/finance/chart/TTTT?period1=pppppppp&period2=qqqqqqqq&interval=1d&events=eeeeeeee
```

* symbol - Ticker (e.g., AAPL, MSFT, etc.)
* period1 - Period1 is the timestamp (POSIX time stamp) of the beginning date
* period2 - Period2 is the timestamp (POSIX time stamp) of the ending date
* events - Event, can be one of 'history', 'div', or 'split'
  * TODO: Add Flag to support Yahoo Finance event types: (history, split, div)


### Command-line Example 

```bash
go-stock-cli -symbol AMZN -period1 1704117600 -period2 1710133200 -interval 1d
```

```bash
go-stock-cli -filename "filename.json" -period1 $PERIOD1 -period2 $PERIOD2 -interval $INTERVAL
```


### Docker Example 

```bash
docker run --rm -v $(pwd):/app/data go:stock-cli ./go-stock-cli -symbol $symbol -period1 $PERIOD1 -period2 $PERIOD2 -interval $INTERVAL; 
```

```bash
docker run --rm -v $(pwd):/app/data go:stock-cli ./go-stock-cli -filename "filename.json" -period1 $PERIOD1 -period2 $PERIOD2 -interval $INTERVAL
```
