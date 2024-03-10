package main

import (
	"flag"
	"fmt"
	"github.com/manguilar22/go-stock-cli/stock"
	"os"
)

var symbol string
var period1 string
var period2 string
var interval string

func init() {
	flag.StringVar(&symbol, "symbol", "Symbol", "Stock Ticker Symbol")
	flag.StringVar(&period1, "period1", "Period1", "Start timestamp UNIX")
	flag.StringVar(&period2, "period2", "Period2", "End Timestamp UNIX")
	flag.StringVar(&interval, "interval", "Interval", "Time Interval")
	flag.Parse()
}

func main() {
	if flag.NArg() != 0 {
		fmt.Println("Provide arguments: [symbol, period1, period2, interval]")
		os.Exit(1)
	}

	fileName := fmt.Sprintf("%s_%s_%s_%s.csv", symbol, period1, period2, interval)
	_ = stock.SaveToCSV(symbol, period1, period2, interval, fileName)

}
