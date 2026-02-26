package http

import (
	"basket-service/internal/adapters/in/http/problems"
	"basket-service/mocks/core/application/usecases/commandsmocks"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_Checkout_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	basketID := uuid.New()

	checkoutCommandHandlerMock := &commandsmocks.CheckoutCommandHandlerMock{}
	checkoutCommandHandlerMock.
		On("Handle", mock.Anything, mock.Anything).
		Return(nil).
		Once()

	s := &Server{
		checkoutCommandHandler: checkoutCommandHandlerMock,
	}

	// Act
	err := s.Checkout(ctx, basketID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestServer_Checkout_CommandCreationError(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	// Передадим невалидный UUID (например, пустой)
	basketID := uuid.Nil

	checkoutCommandHandlerMock := &commandsmocks.CheckoutCommandHandlerMock{}
	checkoutCommandHandlerMock.
		On("Handle", mock.Anything, mock.Anything).
		Return(nil).
		Once()

	s := &Server{
		checkoutCommandHandler: checkoutCommandHandlerMock,
	}

	// Act
	err := s.Checkout(ctx, basketID)

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, problems.ProblemBadRequest)
}
