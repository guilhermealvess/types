package types

import (
	"math"
	"testing"
)

func TestFrom(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected MoneyBRL
	}{
		{
			name:     "zero int",
			input:    0,
			expected: 0,
		},
		{
			name:     "one real as int",
			input:    1,
			expected: Real,
		},
		{
			name:     "negative int",
			input:    -42,
			expected: -42 * Real,
		},
		{
			name:     "int32",
			input:    int32(100),
			expected: 100 * Real,
		},
		{
			name:     "int64",
			input:    int64(100),
			expected: 100 * Real,
		},
		{
			name:     "uint",
			input:    uint(10),
			expected: 10 * Real,
		},
		{
			name:     "uint8",
			input:    uint8(255),
			expected: 255 * Real,
		},
		{
			name:     "uint16",
			input:    uint16(1000),
			expected: 1000 * Real,
		},
		{
			name:     "uint32",
			input:    uint32(50000),
			expected: 50000 * Real,
		},
		{
			name:     "uint64",
			input:    uint64(999),
			expected: 999 * Real,
		},
		{
			name:     "uintptr",
			input:    uintptr(7),
			expected: 7 * Real,
		},
		{
			name:     "one real as float64",
			input:    1.0,
			expected: Real,
		},
		{
			name:     "float64 with four decimal places",
			input:    1.2345,
			expected: 12345,
		},
		{
			name:     "float64 rounds half up",
			input:    1.23456,
			expected: 12346,
		},
		{
			name:     "float32",
			input:    float32(0.5),
			expected: 5000,
		},
		{
			name:     "negative float",
			input:    -3.1415,
			expected: -31415,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got MoneyBRL
			switch v := tt.input.(type) {
			case int:
				got = From(v)
			case int32:
				got = From(v)
			case int64:
				got = From(v)
			case uint:
				got = From(v)
			case uint8:
				got = From(v)
			case uint16:
				got = From(v)
			case uint32:
				got = From(v)
			case uint64:
				got = From(v)
			case uintptr:
				got = From(v)
			case float64:
				got = From(v)
			case float32:
				got = From(v)
			default:
				t.Fatalf("unsupported test input type %T", v)
			}

			if got != tt.expected {
				t.Errorf("From() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestMoneyBRL_Real(t *testing.T) {
	tests := []struct {
		name     string
		money    MoneyBRL
		expected float64
	}{
		{
			name:     "zero",
			money:    0,
			expected: 0,
		},
		{
			name:     "one real",
			money:    Real,
			expected: 1.0,
		},
		{
			name:     "one cent",
			money:    Cent,
			expected: 0.01,
		},
		{
			name:     "four decimal precision",
			money:    12345,
			expected: 1.2345,
		},
		{
			name:     "negative",
			money:    -Real,
			expected: -1.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.money.Real()
			if math.Abs(got-tt.expected) > 1e-9 {
				t.Errorf("Real() = %.10f, want %.10f", got, tt.expected)
			}
		})
	}
}

func TestMoneyBRL_String(t *testing.T) {
	tests := []struct {
		name     string
		money    MoneyBRL
		expected string
	}{
		{
			name:     "zero",
			money:    0,
			expected: "0.00",
		},
		{
			name:     "one real",
			money:    Real,
			expected: "1.00",
		},
		{
			name:     "one cent",
			money:    Cent,
			expected: "0.01",
		},
		{
			name:     "negative",
			money:    -123 * Real,
			expected: "-123.00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.money.String()
			if got != tt.expected {
				t.Errorf("String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestMoneyBRL_Add(t *testing.T) {
	tests := []struct {
		name     string
		a        MoneyBRL
		b        MoneyBRL
		expected MoneyBRL
	}{
		{
			name:     "two positives",
			a:        100 * Real,
			b:        50 * Real,
			expected: 150 * Real,
		},
		{
			name:     "with negative",
			a:        100 * Real,
			b:        -30 * Real,
			expected: 70 * Real,
		},
		{
			name:     "zero",
			a:        0,
			b:        42 * Real,
			expected: 42 * Real,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Add(tt.b)
			if got != tt.expected {
				t.Errorf("Add() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestMoneyBRL_Abs(t *testing.T) {
	tests := []struct {
		name     string
		money    MoneyBRL
		expected MoneyBRL
	}{
		{
			name:     "positive",
			money:    100,
			expected: 100,
		},
		{
			name:     "negative",
			money:    -100,
			expected: 100,
		},
		{
			name:     "zero",
			money:    0,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.money.Abs()
			if got != tt.expected {
				t.Errorf("Abs() = %d, want %d", got, tt.expected)
			}
		})
	}
}

func TestMoneyBRL_Mult(t *testing.T) {
	tests := []struct {
		name     string
		a        MoneyBRL
		b        MoneyBRL
		expected MoneyBRL
	}{
		{
			name:     "one real times two",
			a:        Real,
			b:        2 * Real,
			expected: 2 * Real,
		},
		{
			name:     "two reals times half",
			a:        2 * Real,
			b:        From(0.5),
			expected: Real,
		},
		{
			name:     "zero",
			a:        100 * Real,
			b:        0,
			expected: 0,
		},
		{
			name:     "negative",
			a:        -2 * Real,
			b:        3 * Real,
			expected: -6 * Real,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.a.Mult(tt.b)
			if got != tt.expected {
				t.Errorf("Mult() = %d, want %d", got, tt.expected)
			}
		})
	}
}
