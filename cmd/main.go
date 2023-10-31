package main

import (
	postgres "API/internal/DB"
	server "API/internal/app"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	s := server.New()
	db := postgres.Init()
	if err := postgres.ConnectDB(db); err != nil {
		log.Fatal("Exit")
	}
	if err := postgres.CreateTable(db); err != nil {
		log.Fatal(err)
	}
	if err := s.Start(db); err != nil {
		log.Fatal(err)
	}
	defer postgres.CloseDB(db)
}
