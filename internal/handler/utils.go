package handler

import (
	"fmt"
	"github.com/Di-nis/shortener-url/internal/models"
)

func addBaseURLToShort(baseURL string, urls []models.URL) {
	for idx, url := range urls {
		urls[idx].Short = fmt.Sprintf("%s/%s", baseURL, url.Short)
	}
}
