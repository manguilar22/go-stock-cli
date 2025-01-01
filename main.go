package main

import (
	"flag"
	"fmt"
	"github.com/manguilar22/go-stock-cli/stock"
	"github.com/manguilar22/go-stock-cli/stock/databases"
	"log"
	"os"
	"strconv"
)

var symbol string
var period1 string
var period2 string
var interval string
var filename string
var parallel bool
var dbHost string
var dbPort string
var dbDatabase string
var dbUser string
var dbPassword string

func init() {
	flag.StringVar(&symbol, "symbol", "", "Stock Ticker Symbol")
	flag.StringVar(&period1, "period1", "", "Start timestamp UNIX")
	flag.StringVar(&period2, "period2", "", "End Timestamp UNIX")
	flag.StringVar(&interval, "interval", "", "Time Interval")
	flag.StringVar(&filename, "filename", "", "JSON Filename")
	flag.BoolVar(&parallel, "parallel", false, "Process files in parallel")
	flag.StringVar(&dbHost, "host", "", "Database hostname")
	flag.StringVar(&dbPort, "port", "", "Database port")
	flag.StringVar(&dbDatabase, "database", "", "Database name")
	flag.StringVar(&dbUser, "username", "", "Database login username")
	flag.StringVar(&dbPassword, "password", "", "Database username password")
	flag.Parse()
}

func main() {
	if flag.NArg() != 0 {
		fmt.Println("Provide arguments: [symbol, period1, period2, interval]")
		os.Exit(1)
	}

	if filename != "" && period1 != "" && period2 != "" && interval != "" {
		log.Printf("arguments: -filename %s -period1 %s -period2 %s -interval %s -parallel %t", filename, period1, period2, interval, parallel)
		stock.ProcessFile(filename, period1, period2, interval, parallel)
	}

	if symbol != "" && period1 != "" && period2 != "" && interval != "" {
		log.Printf("arguments: -symbol %s -period1 %s -period2 %s -interval %s", symbol, period1, period2, interval)
		err := stock.SaveToCSV(symbol, period1, period2, interval)
		if err != nil {
			log.Println(err)
		}
	}

	if symbol != "" && dbHost != "" && dbPort != "" && dbDatabase != "" && dbUser != "" && dbPassword != "" {
		port, _ := strconv.Atoi(dbPort)

		config := databases.PostgresConfiguration{
			Host:     dbHost,
			Port:     port,
			DBName:   dbDatabase,
			User:     dbUser,
			Password: dbPassword,
		}

		db := databases.NewDatabase(&config)
		err := db.Connect()

		log.Println("connection status: ", err)

		status, err := db.Exists(symbol)
		log.Printf("does %s exist: status=%t, error=%v", symbol, status, err)

		if status {
			log.Printf("updating records in the %s column: %v", symbol, err)
			err = db.Update(symbol)
			log.Printf("database update: %v", err)
		} else {
			log.Printf("create records for %s: error=%v", symbol, err)
			err = db.Write(symbol)
		}
		db.Close()
	}

}
