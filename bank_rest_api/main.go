package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/meihern/go_learning/api"
	"github.com/meihern/go_learning/storage"
)

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf(err.Error())
	}

	store, err := storage.NewPostgresStore()

	if err != nil {
		log.Fatalf(err.Error())
	}

	if err := store.Init(); err != nil {
		log.Fatalf(err.Error())
	}

	server := api.NewAPIServer(os.Getenv("LISTEN_ADDRESS"), store)

	server.Run()

}
