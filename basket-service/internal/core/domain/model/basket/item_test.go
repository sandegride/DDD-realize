package basket

import (
	"basket-service/internal/core/domain/model/good"
	"basket-service/internal/pkg/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ItemBeCorrectWhenParamsAreCorrectOnCreated(t *testing.T) {
	// Arrange
	coffee := good.Coffee()
	quantity := 1

	// Act
	item, err := NewItem(coffee, quantity)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, item.ID())
	assert.Equal(t, coffee.ID(), item.goodID)
	assert.Equal(t, coffee.Title(), item.title)
	assert.Equal(t, coffee.Description(), item.description)
	assert.Equal(t, coffee.Price(), item.price)
	assert.Equal(t, quantity, item.quantity)
}

func Test_ItemReturnErrorWhenParamsAreNotCorrectOnCreated(t *testing.T) {
	// Arrange
	coffee := good.Coffee()
	tests := map[string]struct {
		good     *good.Good
		quantity int
		expected error
	}{
		"wrong_good": {
			good:     nil,
			quantity: 1,
			expected: errs.NewValueIsRequiredError("good"),
		},
		"wrong_quantity_zero": {
			good:     coffee,
			quantity: 0,
			expected: ErrQuantityIsZeroOrLess,
		},
		"wrong_quantity_less_zero": {
			good:     coffee,
			quantity: -1,
			expected: ErrQuantityIsZeroOrLess,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Act
			_, err := NewItem(test.good, test.quantity)

			// Assert
			if err.Error() != test.expected.Error() {
				t.Errorf("expected %v, got %v", test.expected, err)
			}
		})
	}
}

func Test_ItemCanSetQuantityWhenQuantityIsCorrect(t *testing.T) {
	// Arrange
	coffee := good.Coffee()
	item, _ := NewItem(coffee, 1)

	// Act
	err := item.setQuantity(2)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, item.quantity)
}

func Test_ItemReturnErrorWhenParamsAreNotCorrectOnSetQuantity(t *testing.T) {
	// Arrange
	// Arrange
	coffee := good.Coffee()
	item, _ := NewItem(coffee, 1)

	tests := map[string]struct {
		good     *good.Good
		quantity int
		expected error
	}{
		"wrong_quantity_zero": {
			good:     coffee,
			quantity: 0,
			expected: ErrQuantityIsZeroOrLess,
		},
		"wrong_quantity_less_zero": {
			good:     coffee,
			quantity: -1,
			expected: ErrQuantityIsZeroOrLess,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Act
			err := item.setQuantity(test.quantity)

			// Assert
			if err.Error() != test.expected.Error() {
				t.Errorf("expected %v, got %v", test.expected, err)
			}
		})
	}
}

func Test_ItemShouldCalculateTotal(t *testing.T) {
	// Arrange
	coffee := good.Coffee()
	item, _ := NewItem(coffee, 2)

	// Act
	total := item.getTotal()

	// Assert
	assert.Equal(t, 1000.0, total)
}

func Test_ItemShouldBeEqualWhenIdIsSame(t *testing.T) {
	// Arrange
	coffee := good.Coffee()
	item, _ := NewItem(coffee, 2)

	// Act
	isEqual := item.Equal(item)

	// Assert
	assert.True(t, isEqual)
}

func Test_ItemShouldBeNotEqualWhenIdIsNotSame(t *testing.T) {
	// Arrange
	coffee := good.Coffee()
	item1, _ := NewItem(coffee, 2)
	item2, _ := NewItem(coffee, 2)

	// Act
	isEqual := item1.Equal(item2)

	// Assert
	assert.False(t, isEqual)
}

func Test_ItemShouldBeNotEqualWhenOtherIsNil(t *testing.T) {
	// Arrange
	coffee := good.Coffee()
	item, _ := NewItem(coffee, 2)

	// Act
	isEqual := item.Equal(nil)

	// Assert
	assert.False(t, isEqual)
}
