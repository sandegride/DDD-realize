package courier

import (
	"delivery-service/internal/core/domain/model/kernel"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_CourierShouldBeCorrectWhenParamsAreCorrectOnCreated(t *testing.T) {
	courier, err := NewCourier("Пешеход", 1, kernel.MinLocation())

	assert.NoError(t, err)
	assert.NotNil(t, courier)
	assert.Equal(t, "Пешеход", courier.Name())
	assert.Equal(t, 1, courier.Speed())
	assert.Equal(t, kernel.MinLocation(), courier.Location())
}

func Test_CourierCanCalculateTimeToLocation(t *testing.T) {
	courier, err := NewCourier("Велосипедист", 2, kernel.MinLocation())
	assert.NoError(t, err)
	target, err := kernel.NewLocation(5, 10)
	assert.NoError(t, err)

	time, err := courier.CalculateTimeToLocation(target)

	assert.NoError(t, err)
	assert.Equal(t, 6.5, time)
}
