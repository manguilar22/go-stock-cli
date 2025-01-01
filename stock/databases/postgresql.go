package databases

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"strconv"
)

type PostgresConfiguration struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type Database struct {
	config *PostgresConfiguration
	pool   *pgxpool.Pool
}

func NewDatabase(config *PostgresConfiguration) *Database {
	return &Database{config: config}
}

func (db *Database) Connect() error {
	// Create the connection string
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		db.config.User, db.config.Password, db.config.Host, db.config.Port, db.config.DBName)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return fmt.Errorf("unable to create the connection pool: %w", err)
	}
	db.pool = pool

	err = db.pool.Ping(context.Background())
	if err != nil {
		return fmt.Errorf("unable to establish a connection to the database: %w", err)
	}

	status, err := db.TableExists("stocks")

	if !status {
		log.Println("Create stocks table")
		bytes, _ := os.ReadFile("stock/databases/stocks.sql")
		err = db.CreateTable(string(bytes))
		log.Fatalln("creating database table: ", err)
	}

	log.Println("Connected to the database successfully")
	return nil
}

func (db *Database) Close() {
	db.pool.Close()
	log.Println("Closed connection to database")
}

func (db *Database) CreateTable(schema string) error {
	_, err := db.pool.Exec(context.Background(), schema)

	if err != nil {
		return fmt.Errorf("failed to create database table from schema: %w", err)
	}
	return nil
}

func (db *Database) TableExists(tableName string) (bool, error) {
	var status bool
	query := `SELECT EXISTS(
		SELECT 1
		FROM information_schema.tables 
		WHERE table_schema = $1 
		AND table_name = $2
	);`

	err := db.pool.QueryRow(context.Background(), query, "public", tableName).Scan(&status)
	if err != nil {
		return false, fmt.Errorf("error checking if table exists: %w", err)
	}

	return status, nil

}

func (db *Database) Exists(symbol string) (bool, error) {
	var status bool

	queryString := `SELECT EXISTS(SELECT 1 FROM stocks WHERE symbol = $1)`
	err := db.pool.QueryRow(context.Background(), queryString, symbol).Scan(&status)

	if err != nil {
		return false, err
	}

	return status, nil
}

func (db *Database) Update(symbol string) error {
	filePath := fmt.Sprintf("data/csv/%s.csv", symbol)

	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		log.Fatalf("Failed to update %s: %w", symbol, err)
		return err
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		log.Printf("failed to update %s from the database: %w", symbol, err)
	}

	currentPeriod2 := records[1][2]

	period2CheckSql := fmt.Sprintf("SELECT period2 FROM stocks WHERE symbol = '%s' LIMIT 1", symbol)
	var period2CheckValue string

	err = db.pool.QueryRow(context.Background(), period2CheckSql).Scan(&period2CheckValue)

	if err != nil {
		log.Println("failed to check if period2 is the same: ", period2CheckSql)
	}

	if currentPeriod2 == period2CheckValue {
		return fmt.Errorf("the %s symbol does not need an update", symbol)
	}

	currentSymbol := records[1][0]
	sqlString := fmt.Sprintf("DELETE FROM stocks WHERE symbol = '%s'", currentSymbol)
	log.Println("SQL STATEMENT: ", sqlString)

	deleteOutput, err := db.pool.Exec(context.Background(), sqlString)
	if err != nil {
		log.Println("failed to delete stock symbol: ", symbol)
		return err
	}

	log.Println("cleaned stocks table: ", deleteOutput)

	err = db.Write(symbol)

	if err != nil {
		log.Println("failed to upload stock symbol: ", symbol)
		return err
	}

	return nil
}

func (db *Database) Write(symbol string) error {
	directory, _ := os.Getwd()
	var filePath string = fmt.Sprintf("%s/data/csv/%s.csv", directory, symbol)

	stockFile, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("%s does not exist.", filePath)
	}

	log.Printf("%s symbol exists with %s permissions", symbol, stockFile.Mode().String())
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("failure to read file: %s", file.Name())
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Unable to parse file as CSV: %v\n", err)
	}

	for i, record := range records {
		fmt.Printf("%d = %v\n", i, record)

		// Extract values from the record
		symbol := record[0]
		period1, _ := strconv.Atoi(record[1])
		period2, _ := strconv.Atoi(record[2])
		interval := record[3]
		date := record[4]
		open, _ := strconv.ParseFloat(record[5], 64)
		high, _ := strconv.ParseFloat(record[6], 64)
		low, _ := strconv.ParseFloat(record[7], 64)
		close, _ := strconv.ParseFloat(record[8], 64)
		volume, _ := strconv.ParseFloat(record[9], 64)

		// SQL INSERT statement
		sqlString := fmt.Sprintf(`
			INSERT INTO stocks(symbol, period1, period2, interval, date, open, high, low, close, volume)
			VALUES ('%s', %d, %d, '%s', '%s', %f, %f, %f, %f, %f)
		`,
			symbol,
			period1,
			period2,
			interval,
			date,
			open,
			high,
			low,
			close,
			volume,
		)

		//log.Println(fmt.Sprintf("SQL string: %s", sqlString))
		_, err := db.pool.Exec(context.Background(), sqlString)
		if err != nil {
			log.Printf("Record %d has error: %v", i, err)
		}
	}
	return nil
}
