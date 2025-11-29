package storage

import (
	"github.com/Di-nis/shortener-url/internal/models"
)

// ProducerMemory - структура для записи данных в файл.
type ProducerMemory struct {
	URLs []models.URLBase
}

func NewProducerMemory(urls []models.URLBase) *ProducerMemory {
	return &ProducerMemory{
		URLs: urls,
	}
}

// Write - запись данных в файл.
func (p *ProducerMemory) Write(url models.URLBase) error {
	p.URLs = append(p.URLs, url)
	return nil
}

// Close - закрытие файла.
func (p *ProducerMemory) Close() error {
	return nil
}

// ConsumerMemory - структура для чтения данных из файла.
type ConsumerMemory struct {
	URLs []models.URLBase
}

func NewConsumerMemory(urls []models.URLBase) *ConsumerMemory {
	return &ConsumerMemory{
		URLs: urls,
	}
}

func (c *ConsumerMemory) Load() ([]models.URLBase, error) {
	return c.URLs, nil
}

func (c *ConsumerMemory) Close() error {
	return nil
}
