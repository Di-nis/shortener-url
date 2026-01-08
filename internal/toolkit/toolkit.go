package toolkit

import "fmt"

// addBaseURLToResponse - добавление базового URL к ответу.
func AddBaseURLToResponse(baseURL string, urlShort string) string {
	return fmt.Sprintf("%s/%s", baseURL, urlShort)
}
