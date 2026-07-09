package types

import (
	"database/sql/driver"
	"fmt"
)

type Email string

func NewEmail(value string) (Email, error) {
	e := Email(value)
	if err := validate.Var(e, "email"); err != nil {
		return "", fmt.Errorf("invalid email: %w", err)
	}
	return e, nil
}

func (e Email) String() string {
	return string(e)
}

func (e *Email) Scan(v interface{}) error {
	switch x := v.(type) {
	case string:
		email, err := NewEmail(x)
		if err != nil {
			return err
		}
		*e = email
	case []byte:
		return e.Scan(string(x))
	default:
		return fmt.Errorf("cannot scan %T into Email", v)
	}
	return nil
}

func (e Email) Value() (driver.Value, error) {
	return string(e), nil
}
