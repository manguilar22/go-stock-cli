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
		log.Printf("filename=%s,period1=%s,period2=%s,interal=%s", filename, period1, period2, interval)
		stock.ProcessFile(filename, period1, period2, interval)
	}

	if symbol != "" && period1 != "" && period2 != "" && interval != "" {
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
		log.Printf("does %s exist: %v", symbol, err)

		if status {
			err = db.Update(symbol)
			log.Printf("updating records in the %s column: %v", symbol, err)
		} else {
			err = db.Write(symbol)
			log.Printf("failed to create records for %s: %v", symbol, err)
		}
		db.Close()
	}

}
