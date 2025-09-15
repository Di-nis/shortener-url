package repository

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/Di-nis/shortener-url/internal/models"
)

type Producer struct {
	file   *os.File
	writer *bufio.Writer
}

func NewProducer(filename string) (*Producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &Producer{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

func (p *Producer) WriteURL(url models.URL) error {
	urlTypeTwo := models.URLCopyTwo(url)
	data, err := json.Marshal(&urlTypeTwo)
	if err != nil {
		return err
	}

	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}

	return p.writer.Flush()
}

func (p *Producer) Close() error {
	return p.file.Close()
}

func (p *Producer) SaveToFile(urlData models.URL) error {
	err := p.WriteURL(urlData)
	if err != nil {
		return err
	}

	return nil
}

type Consumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *Consumer) ReadURL() (*models.URL, error) {
	data := c.scanner.Bytes()

	urlsTypeTwo := models.URLCopyTwo{}
	err := json.Unmarshal(data, &urlsTypeTwo)
	if err != nil {
		return nil, err
	}
	urls := models.URL(urlsTypeTwo)

	return &urls, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}

func (c *Consumer) LoadFromFile() ([]models.URL, error) {
	URLArray := make([]models.URL, 0)

	for c.scanner.Scan() {
		urlData, err := c.ReadURL()
		if err != nil {
			return nil, err
		}
		URLArray = append(URLArray, *urlData)
	}
	err := c.Close()
	if err != nil {
		return nil, err
	}
	return URLArray, nil
}
