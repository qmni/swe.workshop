package httpapi

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestCreateProductRequestValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		request createProductRequest
		wantErr bool
	}{
		{
			name: "valid product",
			request: createProductRequest{
				Name:        "Notebook",
				Description: "Workshop product",
				PriceCents:  1299,
			},
			wantErr: false,
		},
		{
			name: "missing name",
			request: createProductRequest{
				Description: "Workshop product",
				PriceCents:  1299,
			},
			wantErr: true,
		},
		{
			name: "invalid price",
			request: createProductRequest{
				Name:       "Notebook",
				PriceCents: 0,
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
