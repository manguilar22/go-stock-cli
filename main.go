package main

import (
	"flag"
	"fmt"
	"github.com/manguilar22/go-stock-cli/stock"
	"log"
	"os"
)

var symbol string
var period1 string
var period2 string
var interval string
var filename string

func init() {
	flag.StringVar(&symbol, "symbol", "", "Stock Ticker Symbol")
	flag.StringVar(&period1, "period1", "", "Start timestamp UNIX")
	flag.StringVar(&period2, "period2", "", "End Timestamp UNIX")
	flag.StringVar(&interval, "interval", "", "Time Interval")
	flag.StringVar(&filename, "filename", "", "JSON Filename")
	flag.Parse()
}

func main() {
	if flag.NArg() != 0 {
		fmt.Println("Provide arguments: [symbol, period1, period2, interval]")
		os.Exit(1)
	}

	if filename != "" && period1 != "" && period2 != "" && interval != "" {
		log.Printf("filename=%s,period1=%s,period2=%s,interal=%s", filename, period1, period2, interval)
		stock.ProcessFile(filename, period1, period2, interval)
	}

	if symbol != "" && period1 != "" && period2 != "" && interval != "" {
		err := stock.SaveToCSV(symbol, period1, period2, interval)
		if err != nil {
			log.Println(err)
		}
	}

}
