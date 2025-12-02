// Package main предоставляет точку входа для запуска
// набора статических анализаторов как единого инструмента. Пакет предназначен
// для сборки бинаря, подобного `go vet`, который запускает сразу множество
// анализаторов из `golang.org/x/tools/go/analysis/passes`, `honnef.co/go/tools`
// (stylecheck) и `honnef.co/go/staticcheck` (staticcheck).
package main

import (
	"github.com/Di-nis/shortener-url/internal/multichecker"
)

func main() {
	multichecker.Run()
}
