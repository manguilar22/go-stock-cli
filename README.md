# Yahoo Finance CLI Tool

## Description 

A command-line utility tool for downloading historical stock data from Yahoo Finance. 

### [Yahoo Finance API](https://cryptocointracker.com/yahoo-finance/yahoo-finance-api)

* [https://query1.finance.yahoo.com/v7/finance/download/{symbol}](https://cryptocointracker.com/yahoo-finance/yahoo-finance-api#d7a67d5b007945f286bf8f4e33e89c86)
  * download stock ticker
  * interval: "1m 2m 5m 15m 30m 60m 90m 1h 1d 5d 1wk 1mo 3mo"
  * period1: start time - UNIX timestamp
  * period2: end time - UNIX timestamp

### Command-line Example 

```bash
go-stock-cli -symbol AMZN -period1 1704117600 -period2 1710133200 -interval 1d
```
