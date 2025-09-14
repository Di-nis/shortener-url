package handler

import (
	"fmt"
)

func addBaseURLToShort(baseURL string, urlShort string) string {
	return fmt.Sprintf("%s/%s", baseURL, urlShort)
	
}
