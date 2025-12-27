package main

import (
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "user=postgres password=govno2 dbname=postgres sslmode=disable"

	storage, err := Connect(connStr)
	if err != nil {
		log.Fatal(err)
	}
	storage.Init()
	StartServer(storage)

}
