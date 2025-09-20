package service

import(
	"github.com/Di-nis/shortener-url/internal/constants"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestService_ShortHash(t *testing.T) {
	tests := []struct {
		name string
		data   string
		length int
		baseURL string
		want   string
	}{
		{
			name: "Тест #1",
			data: "https://practicum.yandex.ru",
			length: constants.HashLength,
			baseURL: "http://localhost:8080",
			want: "bTKNZu94",
		},
		{
			name: "Тест #2",
			data: "https://www.sports.ru",
			length: constants.HashLength,
			baseURL: "http://localhost:8080",
			want: "4BeKySvE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService()
			got := service.ShortHash(tt.data, tt.length, tt.baseURL)

			assert.Equal(t, got, tt.want)
		})
	}
}
