package authn_test

import (
	"os"
	"testing"

	"github.com/Di-nis/shortener-url/internal/authn"
	"github.com/stretchr/testify/assert"
)

func TestGetUserID(t *testing.T) {
	secretKey := os.Getenv("JWT_SECRET")

	tests := []struct {
		name        string
		tokenString string
		secretKey   string
		want        string
	}{
		{
			name:        "TestGetUserID, токен невалидный",
			tokenString: "fsdffvdfrgfxvbxdf",
			secretKey:   secretKey,
			want:        "-1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := authn.GetUserID(tt.tokenString, tt.secretKey)
			assert.Equal(t, tt.want, got, "GetUserID() = %v, want %v", got, tt.want)
		})
	}
}
