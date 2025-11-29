package repository

import (
	"context"

	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/models"
)

// WriteCloser - интерфейс для записи в файл.
type WriteCloser interface {
	Write(models.URLBase) error
	Close() error
}

// ReadCloser - интерфейс для чтения из файла.
type ReadCloser interface {
	Load() ([]models.URLBase, error)
	Close() error
}

// Storage - структура для хранения файлов.
type Storage struct {
	Producer WriteCloser
	Consumer ReadCloser
}

// RepoFileMemory - структура базы данных.
type RepoFileMemory struct {
	URLs    []models.URLBase
	Storage *Storage
}

// NewRepoFileMemory - создание структуры RepoFileMemory.
func NewRepoFileMemory(storage *Storage) *RepoFileMemory {
	return &RepoFileMemory{
		URLs:    make([]models.URLBase, 0),
		Storage: storage,
	}
}

// Ping - проверка соединения с базой данных.
func (repo *RepoFileMemory) Ping(ctx context.Context) error {
	return constants.ErrorMethodNotAllowed
}

// InsertBatch - сохранение нескольких URL в базу данных.
func (repo *RepoFileMemory) InsertBatch(ctx context.Context, urls []models.URLBase) error {
	for _, url := range urls {
		for _, urlDB := range repo.URLs {
			if urlDB.Original == url.Original {
				return constants.ErrorURLAlreadyExist
			}
		}

		repo.URLs = append(repo.URLs, url)

		err := repo.Storage.Producer.Write(url)
		if err != nil {
			return err
		}
	}
	return nil
}

// InsertOrdinary - сохранение ординарного URL в базу данных.
func (repo *RepoFileMemory) InsertOrdinary(ctx context.Context, url models.URLBase) error {
	for _, urlDB := range repo.URLs {
		if urlDB.Original == url.Original {
			return constants.ErrorURLAlreadyExist
		}
	}

	repo.URLs = append(repo.URLs, url)

	err := repo.Storage.Producer.Write(url)
	if err != nil {
		return err
	}
	return nil

}

// SelectOriginal - получение оригинального URL из базы данных.
func (repo *RepoFileMemory) SelectOriginal(ctx context.Context, shortURL string) (string, error) {
	for _, url := range repo.URLs {
		if url.Short == shortURL && url.DeletedFlag {
			return "", constants.ErrorURLAlreadyDeleted
		} else if url.Short == shortURL {
			return url.Original, nil
		}
	}
	return "", constants.ErrorURLNotExist
}

// SelectShort - получение оригинального URL из базы данных.
func (repo *RepoFileMemory) SelectShort(ctx context.Context, originalURL string) (string, error) {
	for _, url := range repo.URLs {
		if url.Original == originalURL {
			return url.Short, nil
		}
	}
	return "", constants.ErrorURLNotExist
}

// SelectAll - получение всех когда-либо сокращенных пользователем URL.
func (repo *RepoFileMemory) SelectAll(ctx context.Context, userID string) ([]models.URLBase, error) {
	var urls []models.URLBase

	for _, url := range repo.URLs {
		if url.UUID == userID {
			urls = append(urls, url)
		}
	}
	return urls, nil
}

// Delete - простановка флага удаления.
func (repo *RepoFileMemory) Delete(ctx context.Context, urls []models.URLBase) error {
	for _, url := range urls {
		for i, urlDB := range repo.URLs {
			if urlDB.Short == url.Short && urlDB.UUID == url.UUID && !urlDB.DeletedFlag {
				repo.URLs[i].Original = ""
				repo.URLs[i].DeletedFlag = true
			}
		}
	}
	return nil
}
