package repository

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
	urlOriginal3 = "https://www.fcdin.com/"
	urlShort3    = "kihjTR8h"
	urlOriginal4 = "https://www.sports.ru/"
	urlShort4    = "jkj7fgk2"

	bodyJSONBatch = `[{"correlation_id":"1","original_url":"https://www.khl.ru/"},{"correlation_id":"2","original_url":"https://www.dynamo.ru/"}]`
	bodyJSON1     = `{"url":"https://www.khl.ru/"}`
	bodyJSON2     = `["lJJpJV7h","kiFL71uv"]`
	bodyText1     = `https://maximum.ru/`
	bodyText2     = `https://www.khl.ru/`

	urlIn1 = models.URLBase{
		UUID:     UUID,
		URLID:    "1",
		Original: urlOriginal1}

	urlIn2 = models.URLBase{
		UUID:     UUID,
		URLID:    "2",
		Original: urlOriginal2}

	urlIn3 = models.URLJSON{
		UUID:     UUID,
		Original: urlOriginal1}

	urlIn4 = models.URLBase{
		UUID:     UUID,
		Original: urlOriginal1}

	urlIn5 = models.URLBase{
		UUID:  UUID,
		Short: urlShort1}

	urlIn6 = models.URLBase{
		UUID:  UUID,
		Short: urlShort2}

	urlsIn1 = []models.URLBase{urlIn1, urlIn2}
	urlsIn2 = []models.URLBase{urlIn5, urlIn6}

	urlOut1 = models.URLBase{
		UUID:        UUID,
		URLID:       "1",
		Original:    urlOriginal1,
		Short:       urlShort1,
		DeletedFlag: false,
	}

	urlOut2 = models.URLBase{
		UUID:        UUID,
		URLID:       "2",
		Original:    urlOriginal2,
		Short:       urlShort2,
		DeletedFlag: false,
	}

	urlOut3 = models.URLBase{
		UUID:        UUID,
		Original:    urlOriginal3,
		Short:       urlShort3,
		DeletedFlag: false,
	}

	urlOut4 = models.URLBase{
		UUID:        UUID,
		Original:    urlOriginal4,
		Short:       urlShort4,
		DeletedFlag: true,
	}

	urlsOut1 = []models.URLBase{urlOut1, urlOut2}
	urlsOut2 = []models.URLBase{
		{
			Original: urlOriginal1,
			Short:    urlShort1,
		},
	}
	urlsOut3 = []models.URLBase{urlOut1, urlOut2, urlOut4}
)
