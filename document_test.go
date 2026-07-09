package types

import (
	"bytes"
	"database/sql/driver"
	"testing"
)

func TestNewDocument(t *testing.T) {
	tests := []struct {
		name        string
		value       string
		expectedDoc Document
		expectedErr bool
	}{
		{
			name:        "valid formatted CPF",
			value:       "018.581.296-13",
			expectedDoc: Document{_type: CPF, value: "01858129613"},
			expectedErr: false,
		},
		{
			name:        "valid raw CPF",
			value:       "01858129613",
			expectedDoc: Document{_type: CPF, value: "01858129613"},
			expectedErr: false,
		},
		{
			name:        "invalid CPF digits",
			value:       "018.581.296-14",
			expectedDoc: Document{},
			expectedErr: true,
		},
		{
			name:        "CPF with all same digits",
			value:       "111.111.111-11",
			expectedDoc: Document{},
			expectedErr: true,
		},
		{
			name:        "valid formatted CNPJ",
			value:       "26.637.142/0001-58",
			expectedDoc: Document{_type: CNPJ, value: "26637142000158"},
			expectedErr: false,
		},
		{
			name:        "valid raw CNPJ",
			value:       "26637142000158",
			expectedDoc: Document{_type: CNPJ, value: "26637142000158"},
			expectedErr: false,
		},
		{
			name:        "CPF with first verifier zero",
			value:       "12345678909",
			expectedDoc: Document{_type: CPF, value: "12345678909"},
			expectedErr: false,
		},
		{
			name:        "CNPJ with first verifier zero",
			value:       "01000100000008",
			expectedDoc: Document{_type: CNPJ, value: "01000100000008"},
			expectedErr: false,
		},
		{
			name:        "CPF with both verifiers zero",
			value:       "00000003700",
			expectedDoc: Document{_type: CPF, value: "00000003700"},
			expectedErr: false,
		},
		{
			name:        "CNPJ with second verifier zero",
			value:       "50547736000900",
			expectedDoc: Document{_type: CNPJ, value: "50547736000900"},
			expectedErr: false,
		},
		{
			name:        "invalid CNPJ digits",
			value:       "26.637.142/0001-59",
			expectedDoc: Document{},
			expectedErr: true,
		},
		{
			name:        "invalid length",
			value:       "123456",
			expectedDoc: Document{},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := NewDocument(tt.value)
			if (err != nil) != tt.expectedErr {
				t.Fatalf("NewDocument(%q) error = %v, wantErr %v", tt.value, err, tt.expectedErr)
			}
			if doc != tt.expectedDoc {
				t.Errorf("NewDocument(%q) = %+v, want %+v", tt.value, doc, tt.expectedDoc)
			}
		})
	}
}

func TestDocumentType_String(t *testing.T) {
	tests := []struct {
		name     string
		dt       DocumentType
		expected string
	}{
		{name: "CPF", dt: CPF, expected: "CPF"},
		{name: "CNPJ", dt: CNPJ, expected: "CNPJ"},
		{name: "unknown", dt: DocumentType(99), expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.dt.String()
			if got != tt.expected {
				t.Errorf("String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestDocumentType_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected DocumentType
		wantErr  bool
	}{
		{name: "int64 CPF", input: int64(0), expected: CPF, wantErr: false},
		{name: "int CNPJ", input: int(1), expected: CNPJ, wantErr: false},
		{name: "string CPF", input: "CPF", expected: CPF, wantErr: false},
		{name: "string CNPJ", input: "CNPJ", expected: CNPJ, wantErr: false},
		{name: "bytes CPF", input: []byte("CPF"), expected: CPF, wantErr: false},
		{name: "invalid string", input: "RG", expected: 0, wantErr: true},
		{name: "invalid type", input: 1.5, expected: 0, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dt DocumentType
			err := dt.Scan(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Scan(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && dt != tt.expected {
				t.Errorf("Scan(%v) = %v, want %v", tt.input, dt, tt.expected)
			}
		})
	}
}

func TestDocumentType_Value(t *testing.T) {
	v, err := CPF.Value()
	if err != nil {
		t.Fatalf("Value() unexpected error: %v", err)
	}
	if v != int64(0) {
		t.Errorf("Value() = %v, want %v", v, int64(0))
	}
}

func TestDocumentType_MarshalText(t *testing.T) {
	b, err := CNPJ.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText() unexpected error: %v", err)
	}
	if !bytes.Equal(b, []byte("CNPJ")) {
		t.Errorf("MarshalText() = %q, want %q", b, "CNPJ")
	}
}

func TestDocumentType_UnmarshalText(t *testing.T) {
	tests := []struct {
		name     string
		text     []byte
		expected DocumentType
		wantErr  bool
	}{
		{name: "CPF", text: []byte("CPF"), expected: CPF, wantErr: false},
		{name: "CNPJ", text: []byte("CNPJ"), expected: CNPJ, wantErr: false},
		{name: "invalid", text: []byte("RG"), expected: 0, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dt DocumentType
			err := dt.UnmarshalText(tt.text)
			if (err != nil) != tt.wantErr {
				t.Fatalf("UnmarshalText(%q) error = %v, wantErr %v", tt.text, err, tt.wantErr)
			}
			if !tt.wantErr && dt != tt.expected {
				t.Errorf("UnmarshalText(%q) = %v, want %v", tt.text, dt, tt.expected)
			}
		})
	}
}

func TestDocument_Accessors(t *testing.T) {
	doc := Document{_type: CPF, value: "01858129613"}

	if got := doc.Number(); got != "01858129613" {
		t.Errorf("Number() = %q, want %q", got, "01858129613")
	}
	if got := doc.Type(); got != CPF {
		t.Errorf("Type() = %v, want %v", got, CPF)
	}
	if got := doc.TypeString(); got != "CPF" {
		t.Errorf("TypeString() = %q, want %q", got, "CPF")
	}
}

func TestDocument_String(t *testing.T) {
	tests := []struct {
		name     string
		doc      Document
		expected string
	}{
		{
			name:     "CPF",
			doc:      Document{_type: CPF, value: "01858129613"},
			expected: "CPF: 018.581.296-13",
		},
		{
			name:     "CNPJ",
			doc:      Document{_type: CNPJ, value: "26637142000158"},
			expected: "CNPJ: 26.637.142/0001-58",
		},
		{
			name:     "empty",
			doc:      Document{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.doc.String()
			if got != tt.expected {
				t.Errorf("String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestDocument_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected Document
		wantErr  bool
	}{
		{
			name:     "string CPF",
			input:    "018.581.296-13",
			expected: Document{_type: CPF, value: "01858129613"},
			wantErr:  false,
		},
		{
			name:     "bytes CNPJ",
			input:    []byte("26.637.142/0001-58"),
			expected: Document{_type: CNPJ, value: "26637142000158"},
			wantErr:  false,
		},
		{
			name:     "invalid document",
			input:    "123456",
			expected: Document{},
			wantErr:  true,
		},
		{
			name:     "invalid type",
			input:    123,
			expected: Document{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Document
			err := d.Scan(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Scan(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr && d != tt.expected {
				t.Errorf("Scan(%v) = %+v, want %+v", tt.input, d, tt.expected)
			}
		})
	}
}

func TestDocument_Value(t *testing.T) {
	doc := Document{_type: CPF, value: "01858129613"}
	v, err := doc.Value()
	if err != nil {
		t.Fatalf("Value() unexpected error: %v", err)
	}
	if v != "01858129613" {
		t.Errorf("Value() = %v, want %v", v, "01858129613")
	}
}

func TestDocument_MarshalText(t *testing.T) {
	doc := Document{_type: CPF, value: "01858129613"}
	b, err := doc.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText() unexpected error: %v", err)
	}
	expected := "CPF: 018.581.296-13"
	if string(b) != expected {
		t.Errorf("MarshalText() = %q, want %q", b, expected)
	}
}

func TestDocument_UnmarshalText(t *testing.T) {
	tests := []struct {
		name     string
		text     []byte
		expected Document
		wantErr  bool
	}{
		{
			name:     "valid CPF",
			text:     []byte("018.581.296-13"),
			expected: Document{_type: CPF, value: "01858129613"},
			wantErr:  false,
		},
		{
			name:     "invalid document",
			text:     []byte("123456"),
			expected: Document{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Document
			err := d.UnmarshalText(tt.text)
			if (err != nil) != tt.wantErr {
				t.Fatalf("UnmarshalText(%q) error = %v, wantErr %v", tt.text, err, tt.wantErr)
			}
			if !tt.wantErr && d != tt.expected {
				t.Errorf("UnmarshalText(%q) = %+v, want %+v", tt.text, d, tt.expected)
			}
		})
	}
}

func TestDocument_Mask(t *testing.T) {
	tests := []struct {
		name     string
		doc      Document
		expected string
	}{
		{
			name:     "CPF",
			doc:      Document{_type: CPF, value: "01858129613"},
			expected: "018.581.296-13",
		},
		{
			name:     "CNPJ",
			doc:      Document{_type: CNPJ, value: "26637142000158"},
			expected: "26.637.142/0001-58",
		},
		{
			name:     "unknown type",
			doc:      Document{_type: DocumentType(99), value: "01858129613"},
			expected: "",
		},
		{
			name:     "CPF wrong length",
			doc:      Document{_type: CPF, value: "123"},
			expected: "",
		},
		{
			name:     "CNPJ wrong length",
			doc:      Document{_type: CNPJ, value: "123"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.doc.Mask()
			if got != tt.expected {
				t.Errorf("Mask() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestDocument_Hide(t *testing.T) {
	tests := []struct {
		name     string
		doc      Document
		expected string
	}{
		{
			name:     "CPF",
			doc:      Document{_type: CPF, value: "01858129613"},
			expected: "***581296**",
		},
		{
			name:     "CNPJ",
			doc:      Document{_type: CNPJ, value: "26637142000158"},
			expected: "**6371420001**",
		},
		{
			name:     "unknown type",
			doc:      Document{_type: DocumentType(99), value: "01858129613"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.doc.Hide()
			if got != tt.expected {
				t.Errorf("Hide() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestDocument_MaskAndHide(t *testing.T) {
	tests := []struct {
		name     string
		doc      Document
		expected string
	}{
		{
			name:     "CPF",
			doc:      Document{_type: CPF, value: "01858129613"},
			expected: "***.581.296-**",
		},
		{
			name:     "CNPJ",
			doc:      Document{_type: CNPJ, value: "26637142000158"},
			expected: "**.637.142/0001-**",
		},
		{
			name:     "CPF wrong length",
			doc:      Document{_type: CPF, value: "123"},
			expected: "",
		},
		{
			name:     "CNPJ wrong length",
			doc:      Document{_type: CNPJ, value: "123"},
			expected: "",
		},
		{
			name:     "unknown type",
			doc:      Document{_type: DocumentType(99), value: "01858129613"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.doc.MaskAndHide()
			if got != tt.expected {
				t.Errorf("MaskAndHide() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestDocument_Hide_WrongLength(t *testing.T) {
	doc := Document{_type: CPF, value: "123"}
	if got := doc.Hide(); got != "" {
		t.Errorf("Hide() = %q, want empty string", got)
	}
}

func TestDocument_MaskAndHide_WrongLength(t *testing.T) {
	doc := Document{_type: CNPJ, value: "123"}
	if got := doc.MaskAndHide(); got != "" {
		t.Errorf("MaskAndHide() = %q, want empty string", got)
	}
}

func TestAllSameDigit(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{name: "empty", input: "", expected: true},
		{name: "all same", input: "11111", expected: true},
		{name: "different", input: "12345", expected: false},
		{name: "single char", input: "a", expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := allSameDigit(tt.input)
			if got != tt.expected {
				t.Errorf("allSameDigit(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIsValidCPF_WrongLength(t *testing.T) {
	if isValidCPF("123") {
		t.Error("isValidCPF(123) = true, want false")
	}
}

func TestIsValidCNPJ_WrongLength(t *testing.T) {
	if isValidCNPJ("123") {
		t.Error("isValidCNPJ(123) = true, want false")
	}
}

func TestDocument_ImplementsDriverValue(t *testing.T) {
	var _ driver.Valuer = Document{}
}

func TestDocument_ImplementsScanner(t *testing.T) {
	var _ interface{ Scan(any) error } = (*Document)(nil)
}
