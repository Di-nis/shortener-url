package service

import (
	"github.com/Di-nis/shortener-url/internal/repository"
)

type URLRepository interface {
	GetURL(string) (string, error)
	CreateURL(string) string
}

type Service struct {
	Repo URLRepository
}

func NewService(repo *repository.Repo) *Service {
	return &Service{
		Repo: repo,
	}
}
