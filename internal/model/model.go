package model

// URL 
type URL struct {
	URLOriginal string
	URLShort string
}

// NewURL
func NewURL(urlOriginal, urlShort string) URL {
	return URL{
		URLOriginal: urlOriginal,
		URLShort: urlShort,
	}
}