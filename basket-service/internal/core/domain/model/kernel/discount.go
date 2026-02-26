package kernel

import (
	"basket-service/internal/pkg/errs"
	"math"
)

type Discount struct {
	value float64

	valid bool
}

func NewDiscount(value float64) (Discount, error) {
	const minDiscount = 0
	const maxDiscount = 1
	if value < minDiscount || value > maxDiscount {
		return Discount{}, errs.NewValueIsOutOfRangeError("value", value, minDiscount, maxDiscount)
	}

	return Discount{value: value, valid: true}, nil
}

func (d Discount) Equal(other Discount) bool {
	return d.value == other.value
}

func (d Discount) IsValid() bool {
	return d.valid
}

func (d Discount) IsMore(other Discount) bool {
	return other.value > d.value
}

func (d Discount) IsMoreOrEqual(other Discount) bool {
	return other.value >= d.value
}

func (d Discount) IsLess(other Discount) bool {
	return other.value < d.value
}

func (d Discount) IsLessOrEqual(other Discount) bool {
	return other.value <= d.value
}

func (d Discount) Value() float64 {
	return d.value
}

func (d Discount) Apply(price float64) (float64, error) {
	if price <= 0 {
		return 0, errs.NewValueIsInvalidError("price")
	}
	res := price * (1 - d.value)
	return math.Round(res*100) / 100, nil
}
