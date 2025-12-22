package repository

import (
	"errors"

	"github.com/Di-nis/shortener-url/internal/models"
)

// Тестовые данные
var (
	UUID = "01KA3YRQCWTNAJEGR5Z30PH6VT"

	url1      = "https://www.khl.ru/"
	urlAlias1 = "lJJpJV7h"
	url2      = "https://www.dynamo.ru/"
	urlAlias2 = "kiFL71uv"
	url3      = "https://www.fcdin.com/"
	urlAlias3 = "kihjTR8h"
	url4      = "https://www.sports.ru/"
	urlAlias4 = "jkj7fgk2"

	testURLFull1 = models.URLBase{
		UUID:        UUID,
		URLID:       "1",
		Original:    url1,
		Short:       urlAlias1,
		DeletedFlag: false,
	}

	testURLShort1 = models.URLBase{
		Original: url1,
		Short:    urlAlias1,
	}

	testURLFull2 = models.URLBase{
		UUID:        UUID,
		URLID:       "2",
		Original:    url2,
		Short:       urlAlias2,
		DeletedFlag: false,
	}

	testURLShort2 = models.URLBase{
		Original: url2,
		Short:    urlAlias2,
	}

	testURLFull3 = models.URLBase{
		UUID:        UUID,
		Original:    url3,
		Short:       urlAlias3,
		DeletedFlag: false,
	}

	testURLFull4 = models.URLBase{
		UUID:        UUID,
		Original:    url4,
		Short:       urlAlias4,
		DeletedFlag: true,
	}

	testURLShort4 = models.URLBase{
		Original: url4,
		Short:    urlAlias4,
	}

	testURLsFull  = []models.URLBase{testURLFull1, testURLFull2, testURLFull4}
	testURLsShort = []models.URLBase{testURLShort1, testURLShort2, testURLShort4}
)

// ошибки
var (
	errDB        = errors.New("db error")
	errDBPrepare = errors.New("prepare failed")
)
