package kernel

import (
	"basket-service/internal/pkg/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DiscountBeCorrectWhenParamsAreCorrectInCreate(t *testing.T) {
	discount, err := NewDiscount(1)

	assert.NoError(t, err)
	assert.NotEmpty(t, discount)
	assert.Equal(t, 1.0, discount.Value())
}

func Test_DiscountReturnErrorWhenParamsAreInCorrectOnCreated(t *testing.T) {
	tests := map[string]struct {
		a        float64
		expected error
	}{
		"-1":  {-1, errs.NewValueIsOutOfRangeError("value", -1, 0, 1)},
		"101": {1.1, errs.NewValueIsOutOfRangeError("value", 1.1, 0, 1)},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := NewDiscount(test.a)

			if err.Error() != test.expected.Error() {
				t.Errorf("expected %v, got %v", test.expected, err)
			}
		})
	}
}

func Test_DiscountShouldCanApplyDiscount(t *testing.T) {
	discount, _ := NewDiscount(0.01)
	price := 100.0

	total, err := discount.Apply(price)
	assert.NoError(t, err)
	assert.Equal(t, 99.0, total)
}
