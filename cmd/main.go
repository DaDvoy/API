package main

import (
	server "API/internal/app"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	//_db := DB{}

	s := server.New()
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

	//_db.db, _db.err = sql.Open("postgres", "user=username password=password host=localhost dbname=mydb sslmode=disable")
	//if _db.err != nil {
	//	log.Fatalf("Error: Unable to connect to database: %v", _db.err)
	//}
	//defer _db.db.Close()

}
