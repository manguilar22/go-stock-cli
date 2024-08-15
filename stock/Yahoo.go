package stock

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type StockData struct {
	Date   string `json:"date"`
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Close  string `json:"close"`
	Volume string `json:"volume"`
}

// Get History of security.
func GetStock(stockSymbol, period1, period2, interval string) ([]StockData, error) {
	stockSymbol = sanitizeStockSymbol(stockSymbol)
	// Construct the URL with the provided stock symbol
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v7/finance/download/%s?period1=%s&period2=%s&interval=%s&events=history&includeAdjustedClose=true",
		stockSymbol,
		period1,
		period2,
		interval)

	log.Println(fmt.Sprintf("symbol=%s, url=%s", stockSymbol, url))

	resp, _ := http.Get(url)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("stockSymbol=%s,statusCode=%d,status=%s,url=%s",
			stockSymbol, resp.StatusCode, resp.Status, url)
	}
	defer resp.Body.Close()

	return parseCSV(resp.Body)
}

func sanitizeStockSymbol(stockSymbol string) string {
	if strings.Contains(stockSymbol, ".") {
		var newString string = strings.ReplaceAll(stockSymbol, ".", "-")
		return newString
	}
	return stockSymbol
}

func doesFolderExist(filepath string) error {
	_, err := os.Stat(filepath)

	if err != nil {
		fmt.Println("data/csv folder does not exist.")
		_ = os.Mkdir(filepath, 0777)
	}

	return err
}

func parseCSV(responseBody io.Reader) ([]StockData, error) {
	var stockData []StockData

	csvReader := csv.NewReader(responseBody)

	_, err := csvReader.Read()
	if err != nil {
		return nil, err
	}

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		// Parse CSV record into StockData struct
		data, err := parseRecord(record)
		if err != nil {
			return nil, err
		}

		stockData = append(stockData, data)
	}

	return stockData, nil
}

func parseRecord(record []string) (StockData, error) {
	var stockData StockData

	stockData.Date = record[0]
	stockData.Open = record[1]
	stockData.High = record[2]
	stockData.Low = record[3]
	stockData.Close = record[4]
	stockData.Volume = record[5]

	return stockData, nil
}
