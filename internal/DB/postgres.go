package postgres

import (
	"database/sql"
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
	db  *sql.DB
	err error
}

type Plant struct {
	ID      string  `json:"id"`
	Product string  `json:"product"`
	Amount  string  `json:"amount"`
	Price   float64 `json:"price"`
}

func Init() *DB {
	return &DB{}
}

func ConnectDB(_db *DB) error {
	_db.db, _db.err = sql.Open("postgres", "user=postgres password=postgres host=localhost dbname=postgres sslmode=disable")
	if _db.err != nil {
		log.Fatalf("Error: Unable to connect to database: %v", _db.err)
		return _db.err
	}
	return nil
	//defer _db.db.Close()
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
