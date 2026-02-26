package kernel

import (
	"basket-service/internal/pkg/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	country   = "Россия"
	city      = "Казань"
	street    = "Красносельская"
	house     = "1"
	apartment = "2"
)

func Test_AddressBeCorrectWhenParamsAreCorrectOnCreated(t *testing.T) {
	address, err := NewAddress(country, city, street, house, apartment)

	assert.NoError(t, err)
	assert.NotEmpty(t, address)
	assert.Equal(t, country, address.Country())
	assert.Equal(t, city, address.City())
	assert.Equal(t, street, address.Street())
	assert.Equal(t, house, address.House())
	assert.Equal(t, apartment, address.Apartment())
}

func Test_AddressReturnErrorWhenParamsAreInCorrectOnCreated(t *testing.T) {
	tests := map[string]struct {
		country   string
		city      string
		street    string
		house     string
		apartment string
		expected  error
	}{
		"wrong country": {
			country:   "",
			city:      city,
			street:    street,
			house:     house,
			apartment: apartment,
			expected:  errs.NewValueIsRequiredError("country"),
		},
		"wrong city": {
			country:   country,
			city:      "",
			street:    street,
			house:     house,
			apartment: apartment,
			expected:  errs.NewValueIsRequiredError("city"),
		},
		"wrong street": {
			country:   country,
			city:      city,
			street:    "",
			house:     house,
			apartment: apartment,
			expected:  errs.NewValueIsRequiredError("street"),
		},
		"wrong house": {
			country:   country,
			city:      city,
			street:    street,
			house:     "",
			apartment: apartment,
			expected:  errs.NewValueIsRequiredError("house"),
		},
		"wrong apartment": {
			country:   country,
			city:      city,
			street:    street,
			house:     house,
			apartment: "",
			expected:  errs.NewValueIsRequiredError("apartment"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := NewAddress(tc.country, tc.city, tc.street, tc.house, tc.apartment)

			if err.Error() != tc.expected.Error() {
				t.Errorf("expected: %v, got: %v", tc.expected, err)
			}
		})
	}
}
