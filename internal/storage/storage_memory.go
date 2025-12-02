package storage

import (
	"github.com/Di-nis/shortener-url/internal/models"
)

// ProducerMemory - структура для записи данных в память.
type ProducerMemory struct {
	URLs []models.URLBase
}

// NewProducerMemory - создание нового ProducerMemory.
func NewProducerMemory(urls []models.URLBase) *ProducerMemory {
	return &ProducerMemory{
		URLs: urls,
	}
}

// Write - запись данных в память.
func (p *ProducerMemory) Write(url models.URLBase) error {
	p.URLs = append(p.URLs, url)
	return nil
}

// Close - закрытие.
func (p *ProducerMemory) Close() error {
	return nil
}

// ConsumerMemory - структура для чтения данных из памяти.
type ConsumerMemory struct {
	URLs []models.URLBase
}

// NewConsumerMemory - создание нового ConsumerMemory.
func NewConsumerMemory(urls []models.URLBase) *ConsumerMemory {
	return &ConsumerMemory{
		URLs: urls,
	}
}

// Read - чтение данных из памяти.
func (c *ConsumerMemory) Load() ([]models.URLBase, error) {
	return c.URLs, nil
}

// Close - закрытие.
func (c *ConsumerMemory) Close() error {
	return nil
}
