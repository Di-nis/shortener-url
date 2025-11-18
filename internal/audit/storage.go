package audit

import (
	"bufio"
	"encoding/json"
	"os"
)

// Producer - структура для записи данных в файл.
type Producer struct {
	file   *os.File
	writer *bufio.Writer
}

// NewProducer - создание нового объекта Producer.
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

// Write - запись данных в файл.
func (p *Producer) Write(audit *Audit) error {
	data, err := json.Marshal(&audit)
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

// Close - закрытие файла.
func (p *Producer) Close() error {
	return p.file.Close()
}
