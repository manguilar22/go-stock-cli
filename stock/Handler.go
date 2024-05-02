package stock

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type StockData struct {
	Date   string `json:"date"`
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Close  string `json:"close"`
	Volume string `json:"volume"`
}

func doesFolderExist(filepath string) error {
	_, err := os.Stat(filepath)

	if err != nil {
		fmt.Println("data/csv folder does not exist.")
		_ = os.Mkdir(filepath, 0777)
	}

	return err
}

func SaveToCSV(stockSymbol, period1, period2, interval, fileName string) error {
	_ = doesFolderExist("data/csv")

	records, err := GetStock(stockSymbol, period1, period2, interval)

	if err != nil {
		return fmt.Errorf("(NYSE:%s) does not exist: error=%s", stockSymbol, err.Error())
	}

	var filepath string = fmt.Sprintf("data/csv/%s", fileName)
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Symbol", "Period1", "Period2", "Interval", "Date", "Open", "High", "Low", "Close", "Volume"}
	err = writer.Write(header)
	if err != nil {
		return err
	}

	for _, record := range records {
		var line []string = []string{
			stockSymbol,
			period1,
			period2,
			interval,
			record.Date,
			record.Open,
			record.High,
			record.Low,
			record.Close,
			record.Volume,
		}
		err := writer.Write(line)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetStock(stockSymbol, period1, period2, interval string) ([]StockData, error) {
	// Construct the URL with the provided stock symbol
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v7/finance/download/%s?period1=%s&period2=%s&interval=%s&events=history&includeAdjustedClose=true",
		stockSymbol,
		period1,
		period2,
		interval)

	log.Println(fmt.Sprintf("stocKSymbol=%s, url=%s", stockSymbol, url))

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("stockSymbol=%s,period1=%s,period2=%s,interval=%s,url=%s,error=%s,status=%s",
			stockSymbol, period1, period2, interval, url, err.Error(), resp.Status)
	}
	defer resp.Body.Close()

	return parseCSV(resp.Body)
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
