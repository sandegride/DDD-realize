package courier

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_StoragePlaceShouldBeCorrectWhenParamsAreCorrectOnCreated(t *testing.T) {
	// Arrange

	// Act
	storagePlace, err := NewStoragePlace("Сумка", 10)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, storagePlace)
	assert.Equal(t, "Сумка", storagePlace.Name())
	assert.Equal(t, 10, storagePlace.TotalVolume())
}

func Test_StoragePlaceShouldReturnErrorWhenParamsAreInCorrectOnCreated(t *testing.T) {
	tests := []struct {
		name        string
		totalVolume int
	}{
		{"", 10},
		{"Сумка", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewStoragePlace(tt.name, tt.totalVolume)
			assert.Error(t, err)
		})
	}
}

func Test_StoragePlaceCanStoreShouldWorkCorrectly(t *testing.T) {
	// Arrange
	sp, err := NewStoragePlace("Box", 10)
	assert.NoError(t, err)

	// Act
	canStore, err := sp.CanStore(5)

	// Assert
	assert.NoError(t, err)
	assert.True(t, canStore)
}

func Test_StoragePlaceCanStoreShouldReturnFalseIfVolumeTooLarge(t *testing.T) {
	// Arrange
	sp, err := NewStoragePlace("Box", 10)
	assert.NoError(t, err)

	// Act
	canStore, err := sp.CanStore(20)

	// Assert
	assert.NoError(t, err)
	assert.False(t, canStore)
}

func Test_StoragePlaceStoreShouldFailOnInvalidVolume(t *testing.T) {
	// Arrange
	sp, err := NewStoragePlace("Box", 10)
	assert.NoError(t, err)

	// Act
	err = sp.Store(uuid.New(), 0)

	// Assert
	assert.Error(t, err)
}

func Test_StoragePlaceStoreShouldFailIfCannotStore(t *testing.T) {
	// Arrange
	sp, err := NewStoragePlace("Box", 10)
	assert.NoError(t, err)
	err = sp.Store(uuid.New(), 10)
	assert.NoError(t, err)

	// Act
	err = sp.Store(uuid.New(), 1)

	// Assert
	assert.ErrorIs(t, err, ErrCannotStoreOrderInThisStoragePlace)
}

func Test_StoragePlaceClearShouldSucceed(t *testing.T) {
	// Arrange
	sp, err := NewStoragePlace("Box", 10)
	assert.NoError(t, err)
	orderID := uuid.New()
	err = sp.Store(orderID, 5)
	assert.NoError(t, err)

	// Act
	err = sp.Clear(orderID)

	// Assert
	assert.NoError(t, err)
	assert.Nil(t, sp.OrderID())
}

func Test_StoragePlaceClearShouldFailIfWrongOrderID(t *testing.T) {
	// Arrange
	sp, err := NewStoragePlace("Box", 10)
	assert.NoError(t, err)
	err = sp.Store(uuid.New(), 5)
	assert.NoError(t, err)

	// Act
	err = sp.Clear(uuid.New())

	// Assert
	assert.ErrorIs(t, err, ErrOrderNotStoredInThisPlace)
}

func Test_StoragePlaceClearShouldFailIfOrderIDNil(t *testing.T) {
	// Arrange
	sp, err := NewStoragePlace("Box", 10)
	assert.NoError(t, err)

	// Act
	err = sp.Clear(uuid.Nil)

	// Assert
	assert.Error(t, err)
}
