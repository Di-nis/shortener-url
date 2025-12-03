package main

import (
	"log"

	"github.com/Di-nis/shortener-url/internal/reset"
)

func main() {
	err := reset.App()
	if err != nil {
		log.Fatal(err)
	}
}
