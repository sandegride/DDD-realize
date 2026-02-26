package basket

import (
	"basket-service/internal/pkg/errs"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DeliveryPeriodShouldReturnCorrectIdAndName(t *testing.T) {
	assert.Equal(t, 1, Morning.ID())
	assert.Equal(t, "morning", Morning.Name())
	assert.Equal(t, 6, Morning.From())
	assert.Equal(t, 12, Morning.To())

	assert.Equal(t, 2, Midday.ID())
	assert.Equal(t, "midday", Midday.Name())
	assert.Equal(t, 12, Midday.From())
	assert.Equal(t, 17, Midday.To())

	assert.Equal(t, 3, Evening.ID())
	assert.Equal(t, "evening", Evening.Name())
	assert.Equal(t, 17, Evening.From())
	assert.Equal(t, 24, Evening.To())

	assert.Equal(t, 4, Night.ID())
	assert.Equal(t, "night", Night.Name())
	assert.Equal(t, 0, Night.From())
	assert.Equal(t, 6, Night.To())
}

func Test_DeliveryPeriodCanBeFoundByName(t *testing.T) {
	tests := map[string]struct {
		name       string
		expectedId int
	}{
		"morning": {
			name:       "morning",
			expectedId: 1,
		},
		"midday": {
			name:       "midday",
			expectedId: 2,
		},
		"evening": {
			name:       "evening",
			expectedId: 3,
		},
		"night": {
			name:       "night",
			expectedId: 4,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			deliveryPeriod, err := GetDeliveryPeriodByName(test.name)

			assert.NoError(t, err)
			if deliveryPeriod.ID() != test.expectedId {
				t.Errorf("expected %v, got %v", test.expectedId, deliveryPeriod.ID())
			}
		})
	}
}

func Test_DeliveryPeriodReturnErrorWhenStatusNotFoundByName(t *testing.T) {
	deliveryPeriod, err := GetDeliveryPeriodByName("wrong")

	assert.True(t, errors.Is(err, errs.ErrObjectNotFound))
	assert.Nil(t, deliveryPeriod)
}

func Test_DeliveryPeriodCanBeFoundById(t *testing.T) {
	tests := map[string]struct {
		id           int
		expectedName string
	}{
		"morning": {
			id:           1,
			expectedName: "morning",
		},
		"midday": {
			id:           2,
			expectedName: "midday",
		},
		"evening": {
			id:           3,
			expectedName: "evening",
		},
		"night": {
			id:           4,
			expectedName: "night",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			deliveryPeriod, err := GetDeliveryPeriodByID(test.id)

			assert.NoError(t, err)
			if deliveryPeriod.Name() != test.expectedName {
				t.Errorf("expected %v, got %v", test.expectedName, deliveryPeriod.Name())
			}
		})
	}
}

func Test_DeliveryPeriodReturnErrorWhenStatusNotFoundById(t *testing.T) {
	deliveryPeriod, err := GetDeliveryPeriodByID(-1)

	assert.True(t, errors.Is(err, errs.ErrObjectNotFound))
	assert.Nil(t, deliveryPeriod)
}

func Test_DeliveryPeriodReturnListOfDeliveryPeriods(t *testing.T) {
	assert.NotEmpty(t, DeliveryPeriods)
}
