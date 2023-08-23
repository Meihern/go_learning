package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/meihern/go_learning/api"
)

func main() {
	
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	s := api.NewAPIServer(os.Getenv("LISTEN_ADDRESS"), nil)

	s.Run()

}
