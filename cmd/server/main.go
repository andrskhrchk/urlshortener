package main

import (
	"log"

	"example.com/internal/server"
	"example.com/internal/storage/postgres"
	_ "github.com/lib/pq"
)

func main() {
	connStr := "user=postgres password=govno2 dbname=postgres sslmode=disable"

	storage, err := postgres.Connect(connStr)
	if err != nil {
		log.Fatal(err)
	}
	storage.Init()

	srv := server.NewServer(storage)
	srv.Start()

}
