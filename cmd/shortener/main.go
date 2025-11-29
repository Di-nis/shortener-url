// Package main — точка входа приложения.
package main

import (
	"log"

	"github.com/Di-nis/shortener-url/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
