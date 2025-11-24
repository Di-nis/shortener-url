package repository

import (
	"context"

	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/Di-nis/shortener-url/internal/models"
)

// WriteCloser - интерфейс для записи в файл.
type WriteCloser interface {
	WriteURL(models.URL) error
	SaveToFile(models.URL) error
	Close() error
}

// ReadCloser - интерфейс для чтения из файла.
type ReadCloser interface {
	ReadURL() (*models.URL, error)
	LoadFromFile() ([]models.URL, error)
	Close() error
}

// Storage - структура для хранения файлов.
type Storage struct {
	Producer WriteCloser
	Consumer ReadCloser
}

// RepoFile - структура базы данных.
type RepoFile struct {
	OriginalAndShortURL []models.URL
	FileStoragePath     string
	Storage             *Storage
}

// NewRepoFile - создание структуры Repo.
func NewRepoFile(fileStoragePath string, storage *Storage) *RepoFile {
	return &RepoFile{
		OriginalAndShortURL: make([]models.URL, 0),
		FileStoragePath:     fileStoragePath,
		Storage:             storage,
	}
}

// Ping - проверка соединения с базой данных.
func (repo *RepoFile) Ping(ctx context.Context) error {
	return constants.ErrorMethodNotAllowed
}

// InsertBatch - сохранение нескольких URL в базу данных.
func (repo *RepoFile) InsertBatch(ctx context.Context, urls []models.URL) error {
	for _, url := range urls {
		for _, urlDB := range repo.OriginalAndShortURL {
			if urlDB.Original == url.Original {
				return constants.ErrorURLAlreadyExist
			}
		}

		repo.OriginalAndShortURL = append(repo.OriginalAndShortURL, url)

		if repo.FileStoragePath != "" {
			err := repo.Storage.Producer.SaveToFile(url)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// InsertOrdinary - сохранение ординарного URL в базу данных.
func (repo *RepoFile) InsertOrdinary(ctx context.Context, url models.URL) error {
	for _, urlDB := range repo.OriginalAndShortURL {
		if urlDB.Original == url.Original {
			return constants.ErrorURLAlreadyExist
		}
	}

	repo.OriginalAndShortURL = append(repo.OriginalAndShortURL, url)

	if repo.FileStoragePath != "" {
		err := repo.Storage.Producer.SaveToFile(url)
		if err != nil {
			return err
		}
	}
	return nil

}

// SelectOriginal - получение оригинального URL из базы данных.
func (repo *RepoFile) SelectOriginal(ctx context.Context, shortURL string) (string, error) {
	for _, url := range repo.OriginalAndShortURL {
		if url.Short == shortURL && url.DeletedFlag {
			return "", constants.ErrorURLAlreadyDeleted
		} else if url.Short == shortURL {
			return url.Original, nil
		}
	}
	return "", constants.ErrorURLNotExist
}

// SelectShort - получение оригинального URL из базы данных.
func (repo *RepoFile) SelectShort(ctx context.Context, originalURL string) (string, error) {
	for _, url := range repo.OriginalAndShortURL {
		if url.Original == originalURL {
			return url.Short, nil
		}
	}
	return "", constants.ErrorURLNotExist
}

// SelectAll - получение всех когда-либо сокращенных пользователем URL.
func (repo *RepoFile) SelectAll(ctx context.Context, userID string) ([]models.URL, error) {
	var urls []models.URL

	for _, url := range repo.OriginalAndShortURL {
		if url.UUID == userID {
			urls = append(urls, url)
		}
	}
	return urls, nil
}

// Delete - простановка флага удаления.
func (repo *RepoFile) Delete(ctx context.Context, urls []models.URL) error {
	for _, url := range urls {
		for i, urlDB := range repo.OriginalAndShortURL {
			if urlDB.Short == url.Short && urlDB.UUID == url.UUID && !urlDB.DeletedFlag {
				repo.OriginalAndShortURL[i].Original = ""
				repo.OriginalAndShortURL[i].DeletedFlag = true
			}
		}
	}
	return nil
}
