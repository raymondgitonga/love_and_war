package main

import (
	"github.com/raymondgitonga/love_and_war/api"
	"log"
)

func main() {
	err := api.NewServer()

	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
