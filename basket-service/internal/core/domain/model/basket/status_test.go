package basket

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_StatusShouldReturnCorrectName(t *testing.T) {
	assert.Equal(t, "", StatusEmpty.String())
	assert.Equal(t, "Created", StatusCreated.String())
	assert.Equal(t, "Confirmed", StatusConfirmed.String())
}

func Test_StatusShouldBeEqualWhenAllPropertiesEqual(t *testing.T) {
	assert.True(t, StatusEmpty.Equal(StatusEmpty))
	assert.True(t, StatusCreated.Equal(StatusCreated))
	assert.True(t, StatusConfirmed.Equal(StatusConfirmed))
}

func Test_StatusShouldBeNotEqualWhenAllPropertiesEqual(t *testing.T) {
	assert.False(t, StatusEmpty.Equal(StatusCreated))
	assert.False(t, StatusCreated.Equal(StatusConfirmed))
	assert.False(t, StatusConfirmed.Equal(StatusEmpty))
}
