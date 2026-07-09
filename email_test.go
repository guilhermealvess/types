package types

import (
	"testing"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected Email
		wantErr  bool
	}{
		{
			name:     "valid email",
			value:    "user@example.com",
			expected: "user@example.com",
			wantErr:  false,
		},
		{
			name:     "valid email with subdomain",
			value:    "user@mail.example.com",
			expected: "user@mail.example.com",
			wantErr:  false,
		},
		{
			name:     "invalid missing at sign",
			value:    "userexample.com",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "invalid missing domain",
			value:    "user@",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "empty email",
			value:    "",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEmail(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewEmail(%q) error = %v, wantErr %v", tt.value, err, tt.wantErr)
			}
			if got != tt.expected {
				t.Errorf("NewEmail(%q) = %q, want %q", tt.value, got, tt.expected)
			}
		})
	}
}

func TestEmail_String(t *testing.T) {
	e := Email("user@example.com")
	if got := e.String(); got != "user@example.com" {
		t.Errorf("String() = %q, want %q", got, "user@example.com")
	}
}

func TestEmail_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected Email
		wantErr  bool
	}{
		{
			name:     "string valid email",
			input:    "user@example.com",
			expected: "user@example.com",
			wantErr:  false,
		},
		{
			name:     "bytes valid email",
			input:    []byte("user@example.com"),
			expected: "user@example.com",
			wantErr:  false,
		},
		{
			name:     "invalid email",
			input:    "invalid",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "invalid type",
			input:    123,
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var e Email
			err := e.Scan(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Scan(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && e != tt.expected {
				t.Errorf("Scan(%v) = %q, want %q", tt.input, e, tt.expected)
			}
		})
	}
}

func TestEmail_Value(t *testing.T) {
	e := Email("user@example.com")
	v, err := e.Value()
	if err != nil {
		t.Fatalf("Value() unexpected error: %v", err)
	}
	if v != "user@example.com" {
		t.Errorf("Value() = %v, want %v", v, "user@example.com")
	}
}
