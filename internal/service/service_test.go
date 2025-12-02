package service

import (
	"testing"

	"github.com/Di-nis/shortener-url/internal/constants"
	"github.com/stretchr/testify/assert"
)

func TestService_ShortHash(t *testing.T) {
	tests := []struct {
		name   string
		data   string
		length int
		want   string
	}{
		{
			name:   "Тест #1",
			data:   "https://practicum.yandex.ru",
			length: constants.HashLength,
			want:   "bTKNZu94",
		},
		{
			name:   "Тест #2",
			data:   "https://www.sports.ru",
			length: constants.HashLength,
			want:   "4BeKySvE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService()
			got := service.ShortHash(tt.data, tt.length)

			assert.Equal(t, got, tt.want)
		})
	}
}

func BenchmarkServiceMethods(b *testing.B) {
	service := NewService()

	b.Run("base62Encode", func(b *testing.B) {
		url, length := "https://practicum.yandex.ru", constants.HashLength
		for i := 0; i < b.N; i++ {
			service.ShortHash(url, length)
		}
	})

	b.Run("ShortHash", func(b *testing.B) {
		url, length := "https://practicum.yandex.ru", constants.HashLength
		for i := 0; i < b.N; i++ {
			service.ShortHash(url, length)
		}
	})
}
