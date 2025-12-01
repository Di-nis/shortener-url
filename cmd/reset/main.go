package main

import (
	"github.com/Di-nis/shortener-url/internal/reset"
	"log"
)

func main() {
	err := reset.App()
	if err != nil {
		log.Fatal(err)
	}
}
