package main

import (
	"os"
)

func main() {
	os.Exit(1) // want "using a direct call to os.Exit in the main function of the main package"
}
