package kernel

import (
	"delivery-service/internal/pkg/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

const minCoordinate = 0
const maxCoordinate = 10

func Test_LocationBeCorrectWhenParamsAreCorrectOnCreated(t *testing.T) {
	location, err := NewLocation(minCoordinate, minCoordinate)

	assert.NoError(t, err)
	assert.NotEmpty(t, location)
	assert.Equal(t, minCoordinate, location.X())
	assert.Equal(t, minCoordinate, location.Y())

	location2, err2 := NewLocation(maxCoordinate, maxCoordinate)

	assert.NoError(t, err2)
	assert.NotEmpty(t, location2)
	assert.Equal(t, maxCoordinate, location2.X())
	assert.Equal(t, maxCoordinate, location2.Y())
}

func Test_LocationReturnErrorWhenParamsAreCorrectOnCreated(t *testing.T) {
	tests := map[string]struct {
		x             int
		y             int
		expectedError error
	}{
		"min x coordinate": {
			x:             minCoordinate - 1,
			y:             minCoordinate,
			expectedError: errs.NewValueIsOutOfRangeError("x", minCoordinate-1, minCoordinate, maxCoordinate),
		},
		"max x coordinate": {
			x:             maxCoordinate + 1,
			y:             maxCoordinate,
			expectedError: errs.NewValueIsOutOfRangeError("x", maxCoordinate+1, minCoordinate, maxCoordinate),
		},
		"min y coordinate": {
			x:             minCoordinate,
			y:             minCoordinate - 1,
			expectedError: errs.NewValueIsOutOfRangeError("y", minCoordinate-1, minCoordinate, maxCoordinate),
		},
		"max y coordinate": {
			x:             maxCoordinate,
			y:             maxCoordinate + 1,
			expectedError: errs.NewValueIsOutOfRangeError("y", maxCoordinate+1, minCoordinate, maxCoordinate),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := NewLocation(tc.x, tc.y)

			if err.Error() != tc.expectedError.Error() {
				t.Errorf("expected: %v, got: %v", tc.expectedError, err)
			}
		})
	}
}

func Test_LocationBeCorrectWhenParamsAreCorrectOnRandomCreated(t *testing.T) {
	location := NewRandomLocation()

	assert.Equal(t, true, location.X() >= minCoordinate && location.Y() >= minCoordinate && location.X() <= maxCoordinate && location.Y() <= maxCoordinate)
}

func Test_LocationReturnCorrectDistanceToTarget(t *testing.T) {
	location, _ := NewLocation(4, 9)
	target, _ := NewLocation(2, 6)

	distance, _ := location.DistanceTo(target)

	assert.Equal(t, 5, distance)
}
