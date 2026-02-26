package kernel

import (
	"basket-service/internal/pkg/errs"
)

type Weight struct {
	value int

	valid bool
}

func NewWeight(value int) (Weight, error) {
	if value <= 0 {
		return Weight{}, errs.NewValueIsInvalidError("value")
	}

	return Weight{value: value, valid: true}, nil
}

func (w Weight) Equal(other Weight) bool {
	return w.value == other.value
}

func (w Weight) IsValid() bool {
	return w.valid
}

func (w Weight) IsMore(other Weight) bool {
	return w.value > other.value
}

func (w Weight) IsMoreOrEqual(other Weight) bool {
	return w.value >= other.value
}

func (w Weight) IsLess(other Weight) bool {
	return w.value < other.value
}

func (w Weight) IsLessOrEqual(other Weight) bool {
	return w.value <= other.value
}

func (w Weight) Value() int {
	return w.value
}
