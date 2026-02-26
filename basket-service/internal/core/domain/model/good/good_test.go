package good

import (
	"basket-service/internal/core/domain/model/kernel"
	"basket-service/internal/pkg/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GoodBeCorrectWhenParamsAreCorrectOnCreated(t *testing.T) {
	id := uuid.MustParse("ec85ceee-f186-4e9c-a4dd-2929e69e586c")
	weight, _ := kernel.NewWeight(6)

	good, err := NewGood(id,
		"Хлеб",
		"Описание хлеба",
		100,
		10,
		weight)

	assert.NoError(t, err)
	assert.NotEmpty(t, good)
	assert.Equal(t, "Хлеб", good.Title())
	assert.Equal(t, "Описание хлеба", good.Description())
	assert.Equal(t, 100.0, good.Price())
	assert.Equal(t, 10, good.Quantity())
	assert.Equal(t, weight, good.Weight())
}

func Test_GoodReturnErrorWhenParamsAreInCorrectOnCreated(t *testing.T) {
	id := uuid.MustParse("ec85ceee-f186-4e9c-a4dd-2929e69e586c")
	title := "Хлеб"
	description := "Описание хлеба"
	price := 100.0
	quantity := 10
	weight, _ := kernel.NewWeight(6)

	tests := map[string]struct {
		id          uuid.UUID
		title       string
		description string
		price       float64
		quantity    int
		weight      kernel.Weight
		expected    error
	}{
		"wrong_id": {
			id:          uuid.Nil,
			title:       title,
			description: description,
			price:       price,
			quantity:    quantity,
			weight:      weight,
			expected:    errs.NewValueIsRequiredError("id"),
		},
		"wrong_title": {
			id:          id,
			title:       "",
			description: description,
			price:       price,
			quantity:    quantity,
			weight:      weight,
			expected:    errs.NewValueIsRequiredError("title"),
		},
		"wrong_description": {
			id:          id,
			title:       title,
			description: "",
			price:       price,
			quantity:    quantity,
			weight:      weight,
			expected:    errs.NewValueIsRequiredError("description"),
		},
		"wrong_price": {
			id:          id,
			title:       title,
			description: description,
			price:       -1.0,
			quantity:    quantity,
			weight:      weight,
			expected:    errs.NewValueIsInvalidError("price"),
		},
		"wrong_quantity": {
			id:          id,
			title:       title,
			description: description,
			price:       price,
			quantity:    -1,
			weight:      weight,
			expected:    errs.NewValueIsInvalidError("quantity"),
		},
		"wrong_weight": {
			id:          id,
			title:       title,
			description: description,
			price:       price,
			quantity:    quantity,
			weight:      kernel.Weight{},
			expected:    errs.NewValueIsRequiredError("weight"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := NewGood(test.id, test.title, test.description, test.price, test.quantity, test.weight)

			if err.Error() != test.expected.Error() {
				t.Errorf("expected %v, got %v", test.expected, err)
			}
		})
	}
}

func Test_GoodCanChangeStocksWhenParamsAreCorrect(t *testing.T) {
	id := uuid.MustParse("ec85ceee-f186-4e9c-a4dd-2929e69e586c")
	weight, _ := kernel.NewWeight(6)
	good, _ := NewGood(id,
		"Хлеб",
		"Описание хлеба",
		100,
		10,
		weight)

	err := good.ChangeStocks(5)

	assert.NoError(t, err)
	assert.Equal(t, 5, good.Quantity())
}

func Test_GoodReturnValueIsRequiredErrorWhenChangeStocksWithZeroQuantity(t *testing.T) {
	good := Coffee()

	err := good.ChangeStocks(-1)

	assert.Error(t, err)
	if err.Error() != errs.NewValueIsInvalidError("quantity").Error() {
		t.Errorf("expected %v, got %v", errs.NewValueIsInvalidError("quantity"), err)
	}
}
