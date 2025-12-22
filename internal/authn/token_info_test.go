package authn

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetClaims(t *testing.T) {
	secretKey := os.Getenv("JWT_SECRET")

	tests := []struct {
		name        string
		tokenString string
		secretKey   string
		wantClaims  *Claims
		wantIsValid bool
	}{
		{
			name:        "TestGetUserID, токен невалидный",
			tokenString: "fsdffvdfrgfxvbxdf",
			secretKey:   secretKey,
			wantClaims:  nil,
			wantIsValid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotClaims, gotIsValid := GetClaims(tt.tokenString, tt.secretKey)

			assert.Equal(t, tt.wantClaims, gotClaims, "GetGetClaims() = %v, want %v", gotClaims, tt.wantClaims)
			assert.Equal(t, tt.wantIsValid, gotIsValid, "GetGetClaims() = %v, want %v", gotIsValid, tt.wantIsValid)
		})
	}
}
