package main

import (
	"github.com/SerjLeo/storage_bot/internal/app"
	"log"
)

func main() {
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
