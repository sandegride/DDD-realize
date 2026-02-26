package kernel

import (
	"basket-service/internal/pkg/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_WeightBeCorrectWhenParamsAreCorrectOnCreate(t *testing.T) {
	weight, err := NewWeight(1)

	assert.NoError(t, err)
	assert.NotEmpty(t, weight)
	assert.Equal(t, 1, weight.Value())
}

func Test_WeightReturnErrorWhenParamsAreInCorrectOnCreated(t *testing.T) {
	tests := map[string]struct {
		a        int
		expected error
	}{
		"0": {
			a:        0,
			expected: errs.NewValueIsInvalidError("value"),
		},
		"-1": {
			a:        -1,
			expected: errs.NewValueIsInvalidError("value"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := NewWeight(tc.a)

			if err.Error() != tc.expected.Error() {
				t.Errorf("expected: %v, got: %v", tc.expected, err)
			}
		})
	}
}
