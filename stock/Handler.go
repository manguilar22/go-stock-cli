package stock

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

func ProcessFile(filename, period1, period2, interval string, parallelStatus bool) {
	fileData, err := os.ReadFile(filename)

	if err != nil {
		log.Printf("Was not able to read JSON file: %s", err.Error())
	}

	var data []interface{}
	_ = json.Unmarshal(fileData, &data)

	if parallelStatus {
		var wg sync.WaitGroup
		for _, record := range data {
			wg.Add(1)

			symbol := record.(map[string]interface{})["symbol"].(string)

			go func() {
				defer wg.Done()
				err := SaveToCSV(symbol, period1, period2, interval)

				if err != nil {
					log.Printf("Error processing file: filename=%s error=%s", filename, err.Error())
				}
			}()
			wg.Wait()
		}
	} else {
		for _, record := range data {
			symbol := record.(map[string]interface{})["symbol"].(string)

			err := SaveToCSV(symbol, period1, period2, interval)
			if err != nil {
				log.Printf("Error processing file: filename=%s error=%s", filename, err.Error())
			}
		}
	}
}

func SaveToCSV(stockSymbol, period1, period2, interval string) error {
	directory, _ := os.Getwd()
	var datadir string = fmt.Sprintf("%s/data/csv", directory)

	_ = DoesFolderExist(datadir)
	records, err := GetStock(stockSymbol, period1, period2, interval)

	if err != nil {
		return fmt.Errorf("(NYSE:%s) does not exist: %s", stockSymbol, err.Error())
	}

	var filepath string = fmt.Sprintf("%s/%s.csv", datadir, stockSymbol)
	log.Printf("Saving %s in %s", stockSymbol, filepath)

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

func DoesFolderExist(filepath string) error {
	_, err := os.Stat(filepath)

	if err != nil {
		_ = os.Mkdir(filepath, 0777)
		return fmt.Errorf("csv folder does not exist.")
	}

	return err
}
