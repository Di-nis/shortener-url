package app

import (
	// "database/sql"
	"fmt"
	"net/http"

	"github.com/Di-nis/shortener-url/internal/config"
	"github.com/Di-nis/shortener-url/internal/handler"
	"github.com/Di-nis/shortener-url/internal/logger"
	"github.com/Di-nis/shortener-url/internal/repository"
	"github.com/Di-nis/shortener-url/internal/service"
	"github.com/Di-nis/shortener-url/internal/usecase"

	"go.uber.org/zap"
)

func initConfigAndLogger() (*config.Config, error) {
	config := &config.Config{}
	config.Parse()

	var err error
	if err = logger.Initialize(config.LogLevel); err != nil {
		return nil, err
	}
	logger.Log.Info("Запуска сервера", zap.String("address", config.ServerAddress))
	return config, nil
}

func initStorageByPostgres(cfg *config.Config) (*repository.RepoPostgres, error) {
	// db, err := sql.Open("pgx", cfg.DataBaseDSN)
	// if err != nil {
	// 	return nil, err
	// }
	// defer db.Close()
	config := repository.NewConfig(cfg.DataBaseDSN)
	repo := repository.NewRepoPostgres(config, cfg.DataBaseDSN)
	return repo, nil
}

func initStorageByFile(cfg *config.Config) (*repository.Repo, error) {
	consumer, err := repository.NewConsumer(cfg.FileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации consumer: %w", err)
	}

	producer, err := repository.NewProducer(cfg.FileStoragePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации producer: %w", err)
	}

	storage := &repository.Storage{
		Consumer: consumer,
		Producer: producer,
	}

	repo := repository.NewRepo(cfg.FileStoragePath, storage)
	repo.URLOriginalAndShort, err = storage.Consumer.LoadFromFile()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки данных из файла-хранилища: %w", err)
	}

	return repo, nil
}

func initStorage(cfg *config.Config) (usecase.URLRepository, error) {
	if cfg.DataBaseDSN == "" {
		return initStorageByFile(cfg)
	}
	return initStorageByPostgres(cfg)
}

func setupRouter(cfg *config.Config, repo usecase.URLRepository, svc *service.Service) http.Handler {
	urlUseCase := usecase.NewURLUseCase(repo, svc)
	controller := handler.NewСontroller(urlUseCase, cfg)
	return controller.CreateRouter()
}
