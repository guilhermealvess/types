package types

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type DocumentType int

func (dt DocumentType) String() string {
	switch dt {
	case CPF:
		return "CPF"
	case CNPJ:
		return "CNPJ"
	}
	return ""
}

func (dt *DocumentType) Scan(v interface{}) error {
	switch x := v.(type) {
	case int64:
		*dt = DocumentType(x)
	case int:
		*dt = DocumentType(x)
	case string:
		switch x {
		case "CPF":
			*dt = CPF
		case "CNPJ":
			*dt = CNPJ
		default:
			return fmt.Errorf("invalid document type: %s", x)
		}
	case []byte:
		return dt.Scan(string(x))
	default:
		return fmt.Errorf("cannot scan %T into DocumentType", v)
	}
	return nil
}

func (dt DocumentType) Value() (driver.Value, error) {
	return int64(dt), nil
}

func (dt DocumentType) MarshalText() ([]byte, error) {
	return []byte(dt.String()), nil
}

func (dt *DocumentType) UnmarshalText(text []byte) error {
	switch string(text) {
	case "CPF":
		*dt = CPF
	case "CNPJ":
		*dt = CNPJ
	default:
		return fmt.Errorf("invalid document type: %s", text)
	}
	return nil
}

const (
	CPF DocumentType = iota
	CNPJ
)

type Document struct {
	_type DocumentType
	value string
}

func NewDocument(value string) (Document, error) {
	cleaned := cleanDocument(value)
	switch len(cleaned) {
	case 11:
		if !isValidCPF(cleaned) {
			return Document{}, fmt.Errorf("invalid CPF: %s", value)
		}
		return Document{_type: CPF, value: cleaned}, nil
	case 14:
		if !isValidCNPJ(cleaned) {
			return Document{}, fmt.Errorf("invalid CNPJ: %s", value)
		}
		return Document{_type: CNPJ, value: cleaned}, nil
	default:
		return Document{}, fmt.Errorf("invalid document length: %s", value)
	}
}

func (d Document) Number() string {
	return d.value
}

func (d Document) Type() DocumentType {
	return d._type
}

func (d Document) TypeString() string {
	return d.Type().String()
}

func (d Document) String() string {
	if d.value == "" {
		return ""
	}
	return d.TypeString() + ": " + d.Mask()
}

func (d *Document) Scan(v any) error {
	switch x := v.(type) {
	case string:
		doc, err := NewDocument(x)
		if err != nil {
			return err
		}
		*d = doc
	case []byte:
		return d.Scan(string(x))
	default:
		return fmt.Errorf("cannot scan %T into Document", v)
	}
	return nil
}

func (d Document) Value() (driver.Value, error) {
	return d.value, nil
}

func (d Document) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *Document) UnmarshalText(text []byte) error {
	doc, err := NewDocument(string(text))
	if err != nil {
		return err
	}
	*d = doc
	return nil
}

// Mask returns the document formatted for display.
// CPF: XXX.XXX.XXX-XX (e.g. 018.581.296-13).
// CNPJ: XX.XXX.XXX/XXXX-XX (e.g. 26.637.142/0001-58).
// Returns an empty string when the document type is unknown.
func (d *Document) Mask() string {
	switch d._type {
	case CPF:
		if len(d.value) != 11 {
			return ""
		}
		return fmt.Sprintf("%s.%s.%s-%s", d.value[0:3], d.value[3:6], d.value[6:9], d.value[9:11])
	case CNPJ:
		if len(d.value) != 14 {
			return ""
		}
		return fmt.Sprintf("%s.%s.%s/%s-%s", d.value[0:2], d.value[2:5], d.value[5:8], d.value[8:12], d.value[12:14])
	}
	return ""
}

// Hide returns the document number masked and hidden.
// CNPJ: **XXXX** (e.g. **2778470001**).
// CPF: ***XXX** (e.g. ***754733**).
// Returns an empty string when the document type is unknown.
func (d Document) Hide() string {
	switch d._type {
	case CPF:
		if len(d.value) != 11 {
			return ""
		}
		return "***" + d.value[3:9] + "**"
	case CNPJ:
		if len(d.value) != 14 {
			return ""
		}
		return "**" + d.value[2:12] + "**"
	}
	return ""
}

// MaskAndHide returns the document formatted for display and hidden.
// CNPJ: **XXXX** (e.g. **2778470001**).
// CPF: ***XXX** (e.g. ***754733**).
// Returns an empty string when the document type is unknown.
func (d *Document) MaskAndHide() string {
	switch d._type {
	case CPF:
		hidden := d.Hide()
		if hidden == "" {
			return ""
		}
		return fmt.Sprintf("%s.%s.%s-%s", hidden[0:3], hidden[3:6], hidden[6:9], hidden[9:11])
	case CNPJ:
		hidden := d.Hide()
		if hidden == "" {
			return ""
		}
		return fmt.Sprintf("%s.%s.%s/%s-%s", hidden[0:2], hidden[2:5], hidden[5:8], hidden[8:12], hidden[12:14])
	}
	return ""
}

func cleanDocument(value string) string {
	var b strings.Builder
	b.Grow(len(value))
	for i := 0; i < len(value); i++ {
		c := value[i]
		if c >= '0' && c <= '9' {
			b.WriteByte(c)
		}
	}
	return b.String()
}

func allSameDigit(s string) bool {
	if len(s) == 0 {
		return true
	}
	first := s[0]
	for i := 1; i < len(s); i++ {
		if s[i] != first {
			return false
		}
	}
	return true
}

func isValidCPF(cpf string) bool {
	if len(cpf) != 11 {
		return false
	}
	if allSameDigit(cpf) {
		return false
	}

	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(cpf[i]-'0') * (10 - i)
	}
	rem := sum % 11
	d1 := 0
	if rem >= 2 {
		d1 = 11 - rem
	}
	if int(cpf[9]-'0') != d1 {
		return false
	}

	sum = 0
	for i := 0; i < 10; i++ {
		sum += int(cpf[i]-'0') * (11 - i)
	}
	rem = sum % 11
	d2 := 0
	if rem >= 2 {
		d2 = 11 - rem
	}
	return int(cpf[10]-'0') == d2
}

func isValidCNPJ(cnpj string) bool {
	if len(cnpj) != 14 {
		return false
	}
	if allSameDigit(cnpj) {
		return false
	}

	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	sum := 0
	for i := 0; i < 12; i++ {
		sum += int(cnpj[i]-'0') * weights1[i]
	}
	rem := sum % 11
	d1 := 0
	if rem >= 2 {
		d1 = 11 - rem
	}
	if int(cnpj[12]-'0') != d1 {
		return false
	}

	sum = 0
	for i := 0; i < 13; i++ {
		sum += int(cnpj[i]-'0') * weights2[i]
	}
	rem = sum % 11
	d2 := 0
	if rem >= 2 {
		d2 = 11 - rem
	}
	return int(cnpj[13]-'0') == d2
}
