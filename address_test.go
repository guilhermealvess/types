package types

import (
	"testing"
)

func TestAddress_String(t *testing.T) {
	tests := []struct {
		name     string
		address  Address
		expected string
	}{
		{
			name: "without complement",
			address: Address{
				Street:       "Av. Paulista",
				Number:       "1000",
				Neighborhood: "Bela Vista",
				City:         "São Paulo",
				State:        "SP",
				Country:      "Brasil",
				ZipCode:      "01310-100",
			},
			expected: "Av. Paulista, 1000 - Bela Vista, São Paulo, SP, Brasil, 01310-100",
		},
		{
			name: "with complement",
			address: Address{
				Street:       "Av. Paulista",
				Number:       "1000",
				Neighborhood: "Bela Vista",
				City:         "São Paulo",
				State:        "SP",
				Country:      "Brasil",
				ZipCode:      "01310-100",
				Complement:   strPtr("Apto 42"),
			},
			expected: "Av. Paulista, 1000 - Bela Vista, São Paulo, SP, Brasil, 01310-100, Apto 42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.address.String()
			if got != tt.expected {
				t.Errorf("String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestAddress_Ok(t *testing.T) {
	tests := []struct {
		name    string
		address Address
		wantErr bool
	}{
		{
			name: "valid address",
			address: Address{
				Street:       "Av. Paulista",
				Number:       "1000",
				Neighborhood: "Bela Vista",
				City:         "São Paulo",
				State:        "SP",
				Country:      "Brasil",
				ZipCode:      "01310-100",
			},
			wantErr: false,
		},
		{
			name: "missing street",
			address: Address{
				Number:       "1000",
				Neighborhood: "Bela Vista",
				City:         "São Paulo",
				State:        "SP",
				Country:      "Brasil",
				ZipCode:      "01310-100",
			},
			wantErr: true,
		},
		{
			name: "missing number",
			address: Address{
				Street:       "Av. Paulista",
				Neighborhood: "Bela Vista",
				City:         "São Paulo",
				State:        "SP",
				Country:      "Brasil",
				ZipCode:      "01310-100",
			},
			wantErr: true,
		},
		{
			name: "missing neighborhood",
			address: Address{
				Street:  "Av. Paulista",
				Number:  "1000",
				City:    "São Paulo",
				State:   "SP",
				Country: "Brasil",
				ZipCode: "01310-100",
			},
			wantErr: true,
		},
		{
			name: "missing city",
			address: Address{
				Street:       "Av. Paulista",
				Number:       "1000",
				Neighborhood: "Bela Vista",
				State:        "SP",
				Country:      "Brasil",
				ZipCode:      "01310-100",
			},
			wantErr: true,
		},
		{
			name: "missing state",
			address: Address{
				Street:       "Av. Paulista",
				Number:       "1000",
				Neighborhood: "Bela Vista",
				City:         "São Paulo",
				Country:      "Brasil",
				ZipCode:      "01310-100",
			},
			wantErr: true,
		},
		{
			name: "missing country",
			address: Address{
				Street:       "Av. Paulista",
				Number:       "1000",
				Neighborhood: "Bela Vista",
				City:         "São Paulo",
				State:        "SP",
				ZipCode:      "01310-100",
			},
			wantErr: true,
		},
		{
			name: "missing zip code",
			address: Address{
				Street:       "Av. Paulista",
				Number:       "1000",
				Neighborhood: "Bela Vista",
				City:         "São Paulo",
				State:        "SP",
				Country:      "Brasil",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.address.Ok()
			if (err != nil) != tt.wantErr {
				t.Errorf("Ok() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}
