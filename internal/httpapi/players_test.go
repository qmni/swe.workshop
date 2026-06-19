package httpapi

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestCreatePlayerRequestValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		request createPlayerRequest
		wantErr bool
	}{
		{
			name: "valid player",
			request: createPlayerRequest{
				Username:    "testplayer",
				Email:       "testplayer@example.com",
				Level:       10,
				Experience:  500,
				PlayerClass: "MAGE",
			},
			wantErr: false,
		},
		{
			name: "missing username",
			request: createPlayerRequest{
				Email:       "testplayer@example.com",
				PlayerClass: "MAGE",
			},
			wantErr: true,
		},
		{
			name: "invalid email",
			request: createPlayerRequest{
				Username:    "testplayer",
				Email:       "not-an-email",
				PlayerClass: "MAGE",
			},
			wantErr: true,
		},
		{
			name: "invalid player class",
			request: createPlayerRequest{
				Username:    "testplayer",
				Email:       "testplayer@example.com",
				PlayerClass: "PALADIN",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.request)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validate.Struct() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
