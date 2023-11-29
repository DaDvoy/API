package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type DB struct {
	db *sql.DB
}

type Plant struct {
	ID      string  `json:"id"`
	Product string  `json:"product"`
	Amount  string  `json:"amount"`
	Price   float64 `json:"price"`
}

func New() (*DB, error) {
	db, err := sql.Open("postgres", "user=postgres password=postgres host=localhost dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatalf("Error: Unable to connect to database: %v", err)
		return nil, fmt.Errorf("$s: %w")
	}
	return &DB{db: db}, nil
}

func CloseDB(_db *DB) error {
	if err := _db.db.Close(); err != nil {
		log.Fatal("Error: Unable to close database:", err)
		return err
	}
	return nil
}

func CreateTable(_db *DB) error {
	rows, err := _db.db.Query("CREATE TABLE IF NOT EXISTS plants(id INT, product VARCHAR, amount INT, price FLOAT)")
	if err != nil {
		log.Fatalf("Error: Unable to create table: %v", err)
		return err
	}
	defer rows.Close()
	return nil
}

func ExecuteQuery(_db *DB, req string, ch chan *Plant) error {
	rows, err := _db.db.Query(req)
	if err != nil {
		log.Fatalf("Error: Unable to execute query: %v", err)
		return err
	}
	if ch != nil {
		var res Plant
		for rows.Next() {
			rows.Scan(&res.ID, &res.Product, &res.Amount, &res.Price)
			ch <- &res
		}
		close(ch)
		defer rows.Close()
	}
	return nil
}
