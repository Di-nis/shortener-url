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

	urlTestData1 = models.URLBase{
		UUID:        UUID,
		URLID:       "1",
		Original:    urlOriginal1,
		Short:       urlShort1,
		DeletedFlag: false,
	}

	urlTestData2 = models.URLBase{
		UUID:        UUID,
		URLID:       "2",
		Original:    urlOriginal2,
		Short:       urlShort2,
		DeletedFlag: false,
	}

	urlTestData3 = models.URLBase{
		UUID:        UUID,
		Original:    urlOriginal3,
		Short:       urlShort3,
		DeletedFlag: false,
	}

	urlTestData4 = models.URLBase{
		UUID:        UUID,
		Original:    urlOriginal4,
		Short:       urlShort4,
		DeletedFlag: true,
	}

	urlsTestData = []models.URLBase{urlTestData1, urlTestData2, urlTestData4}
)
