package handler

import (
	"github.com/Di-nis/shortener-url/internal/models"
)

// Тестовые данные
var (
	UUID = "01KA3YRQCWTNAJEGR5Z30PH6VT"

	urlOriginal1 = "https://www.khl.ru/"
	urlShort1    = "lJJpJV7h"
	urlOriginal2 = "https://www.dynamo.ru/"
	urlShort2    = "kiFL71uv"

	bodyJSONBatch = `[{"correlation_id":"1","original_url":"https://www.khl.ru/"},{"correlation_id":"2","original_url":"https://www.dynamo.ru/"}]`
	bodyJSON1     = `{"url":"https://www.khl.ru/"}`
	bodyJSON2     = `["lJJpJV7h","kiFL71uv"]`
	bodyText1     = `https://maximum.ru/`
	bodyText2     = `https://www.khl.ru/`

	urlIn1 = models.URL{
		UUID:     UUID,
		URLID:    "1",
		Original: urlOriginal1}

	urlIn2 = models.URL{
		UUID:     UUID,
		URLID:    "2",
		Original: urlOriginal2}

	urlIn3 = models.URLCopyOne{
		UUID:     UUID,
		Original: urlOriginal1}

	urlIn4 = models.URL{
		UUID:     UUID,
		Original: urlOriginal1}

	urlIn5 = models.URL{
		UUID:  UUID,
		Short: urlShort1}

	urlIn6 = models.URL{
		UUID:  UUID,
		Short: urlShort2}

	urlsIn1 = []models.URL{urlIn1, urlIn2}
	urlsIn2 = []models.URL{urlIn5, urlIn6}

	urlOut1 = models.URL{
		UUID:        UUID,
		URLID:       "1",
		Original:    urlOriginal1,
		Short:       urlShort1,
		DeletedFlag: false,
	}

	urlOut2 = models.URL{
		UUID:        UUID,
		URLID:       "2",
		Original:    urlOriginal2,
		Short:       urlShort2,
		DeletedFlag: false,
	}

	urlOut3 = models.URL{
		UUID:        UUID,
		Original:    urlOriginal1,
		Short:       urlShort1,
		DeletedFlag: false,
	}

	urlOut4 = models.URL{
		UUID:        UUID,
		Original:    urlOriginal1,
		Short:       urlShort1,
		DeletedFlag: false,
	}

	urlsOut1 = []models.URL{urlOut1, urlOut2}
	urlsOut2 = []models.URL{
		{
			Original: urlOriginal1,
			Short:    urlShort1,
		}}
)
