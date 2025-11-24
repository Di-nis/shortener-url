// Package main — точка входа приложения.
package main

import (
	"github.com/Di-nis/shortener-url/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
