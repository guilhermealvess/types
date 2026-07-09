package types

import (
	"fmt"
	"math"
)

type MoneyBRL int64

const (
	Unit MoneyBRL = 1
	Cent          = 100 * Unit
	Real          = 100 * Cent
)

type Numeric interface {
	int | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64
}

func From[T Numeric](v T) MoneyBRL {
	switch x := any(v).(type) {
	case int:
		return MoneyBRL(x) * Real
	case int32:
		return MoneyBRL(x) * Real
	case int64:
		return MoneyBRL(x) * Real
	case uint:
		return MoneyBRL(x) * Real
	case uint8:
		return MoneyBRL(x) * Real
	case uint16:
		return MoneyBRL(x) * Real
	case uint32:
		return MoneyBRL(x) * Real
	case uint64:
		return MoneyBRL(x) * Real
	case uintptr:
		return MoneyBRL(x) * Real
	case float32:
		return MoneyBRL(math.Round(float64(x) * float64(Real)))
	case float64:
		return MoneyBRL(math.Round(x * float64(Real)))
	}
	return 0
}

// .4 precision
func (m MoneyBRL) Real() float64 {
	return float64(m) / float64(Real)
}

func (m MoneyBRL) String() string {
	return fmt.Sprintf("%.2f", m.Real())
}

func (m MoneyBRL) Add(m2 MoneyBRL) MoneyBRL {
	return m + m2
}

func (m MoneyBRL) Abs() MoneyBRL {
	if m < 0 {
		return m * -1
	}
	return m
}

func (m MoneyBRL) Mult(m2 MoneyBRL) MoneyBRL {
	return MoneyBRL(math.Round(float64(m) * float64(m2) / float64(Real)))
}
