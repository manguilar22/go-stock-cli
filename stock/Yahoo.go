package stock

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type StockData struct {
	Date          string `json:"date"`
	Open          string `json:"open"`
	High          string `json:"high"`
	Low           string `json:"low"`
	Close         string `json:"close"`
	Volume        string `json:"volume"`
	AdjustedClose string `json:"adjustedClose"`
}

// Get History of security.
func GetStock(stockSymbol, period1, period2, interval string) ([]StockData, error) {
	stockSymbol = sanitizeStockSymbol(stockSymbol)
	// Construct the URL with the provided stock symbol
	url := fmt.Sprintf("https://query2.finance.yahoo.com/v8/finance/chart/%s?period1=%s&period2=%s&interval=%s&events=history&includeAdjustedClose=true",
		stockSymbol,
		period1,
		period2,
		interval)

	log.Println(fmt.Sprintf("symbol=%s, url=%s", stockSymbol, url))

	resp, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("unable to perform GET request: %+v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("stockSymbol=%s,statusCode=%d,status=%s,url=%s",
			stockSymbol, resp.StatusCode, resp.Status, url)
	}
	defer resp.Body.Close()

	return parseResponse(resp.Body)
}

func parseResponse(responseBody io.ReadCloser) ([]StockData, error) {
	var response History
	content, err := io.ReadAll(responseBody)

	if err != nil {
		log.Printf("unable convert response content to bytes: %+v", err)
		return nil, err
	}

	err = json.Unmarshal(content, &response)
	if err != nil {
		log.Printf("not able to unmarshall response: %+v", err)
		return nil, err
	}

	if response.Error.Code != "" && response.Error.Description != "" {
		log.Printf("error: code=%s description=%s", response.Error.Code, response.Error.Description)
		return nil, err
	}

	var timestamps []int64 = response.Chart.Result[0].Timestamp
	adjclose := response.Chart.Result[0].Indicators.Adjclose[0]
	quote := response.Chart.Result[0].Indicators.Quote[0]

	var data []StockData = make([]StockData, 0)
	for i := 0; i < len(timestamps); i++ {
		timestamp := timestamps[i]
		unixTimeUTC := time.Unix(timestamp, 0)
		var date string = unixTimeUTC.Format(time.DateOnly)

		var open string = strconv.FormatFloat(quote.Open[i], 'f', -1, 64)
		var closeP string = strconv.FormatFloat(quote.Close[i], 'f', -1, 64)
		var high string = strconv.FormatFloat(quote.High[i], 'f', -1, 64)
		var low string = strconv.FormatFloat(quote.Low[i], 'f', -1, 64)
		var volume string = strconv.FormatFloat(quote.Volume[i], 'f', -1, 64)
		var adjustedClose string = strconv.FormatFloat(adjclose.Adjclose[i], 'f', -1, 64)

		var record StockData = StockData{
			Date:          date,
			Open:          open,
			Close:         closeP,
			High:          high,
			Low:           low,
			Volume:        volume,
			AdjustedClose: adjustedClose,
		}
		data = append(data, record)
	}
	return data, nil
}

func sanitizeStockSymbol(stockSymbol string) string {
	if strings.Contains(stockSymbol, ".") {
		var newString string = strings.ReplaceAll(stockSymbol, ".", "-")
		return newString
	}
	return stockSymbol
}
