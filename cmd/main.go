package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/zeze322/todo/api"
	"github.com/zeze322/todo/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load env file")
	}

	var (
		port        = os.Getenv("TODO_PORT")
		password    = os.Getenv("TODO_PASSWORD")
		storagePath = os.Getenv("TODO_DBFILE")
	)

	store, err := db.NewStorage(storagePath)
	if err != nil {
		log.Println("db error", err)
	}

	defer store.Close()

	s := api.NewServer(port, password, store)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
