package main

import (
	"github.com/meihern/go_learning/api"
)

func main() {
	
	listenAddr := ":3000"

	s := api.NewAPIServer(listenAddr)

	s.Run()

}
