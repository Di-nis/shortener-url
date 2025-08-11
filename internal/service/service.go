package service

import "crypto/sha256"

var base62Alphabet = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") // URLRepository - интерфейс для базы данных.

// Service - структура сервиса по созданию уникального короткого url.
type Service struct{}

// NewService - создание структуры Service.
func NewService() *Service {
	return &Service{}
}

// base62Encode -пПреобразование числа в Base62.
func (service *Service) base62Encode(num uint64) string {
	if num == 0 {
		return string(base62Alphabet[0])
	}
	var encoded []byte
	for num > 0 {
		rem := num % 62
		num /= 62

		encoded = append([]byte{base62Alphabet[rem]}, encoded...)
	}
	return string(encoded)
}

// ShortHash - создание хэш на основе входных данных.
func (service *Service) ShortHash(data string, length int) string {
	var num uint64
	hash := sha256.Sum256([]byte(data))
	for i := range 8 {
		num = (num << 8) | uint64(hash[i])
	}

	b62 := service.base62Encode(num)
	if len(b62) > length {
		return b62[:length]
	}
	return b62
}
