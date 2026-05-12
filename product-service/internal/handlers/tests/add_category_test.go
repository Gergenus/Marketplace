package tests

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gergenus/commerce/product-service/internal/handlers"
	"github.com/Gergenus/commerce/product-service/internal/models"
	"github.com/Gergenus/commerce/product-service/internal/service"
	"github.com/Gergenus/commerce/product-service/internal/service/mocks"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAddCategory(t *testing.T) {

	tests := []struct {
		name                string
		inputBody           string
		inputCategory       models.Category
		mockBehavior        func(ctx context.Context, category string, tst *testing.T) *mocks.MockServiceInterface
		expectedStatusCode  int
		expectedResposeBody string
	}{
		{
			name:          "OK",
			inputBody:     `{"category":"toys"}`,
			inputCategory: models.Category{Category: "toys"},
			mockBehavior: func(ctx context.Context, category string, tst *testing.T) *mocks.MockServiceInterface {
				mock := mocks.NewMockServiceInterface(t)
				mock.EXPECT().AddCategory(ctx, category).Return(1, nil)
				return mock
			},
			expectedStatusCode:  http.StatusOK,
			expectedResposeBody: `{"id":1}` + "\n",
		},
		{
			name:          "Already exists",
			inputBody:     `{"category":"toys"}`,
			inputCategory: models.Category{Category: "toys"},
			mockBehavior: func(ctx context.Context, category string, tst *testing.T) *mocks.MockServiceInterface {
				mock := mocks.NewMockServiceInterface(t)
				mock.EXPECT().AddCategory(ctx, category).Return(-1, service.ErrCategoryAlreadyExists)
				return mock
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedResposeBody: `{"error":"category already exists"}` + "\n",
		},
		{
			name:          "Internal error",
			inputBody:     `{"category":"toys"}`,
			inputCategory: models.Category{Category: "toys"},
			mockBehavior: func(ctx context.Context, category string, tst *testing.T) *mocks.MockServiceInterface {
				mock := mocks.NewMockServiceInterface(t)
				mock.EXPECT().AddCategory(ctx, category).Return(-1, errors.New("internal"))
				return mock
			},
			expectedStatusCode:  http.StatusBadRequest,
			expectedResposeBody: `{"error":"Internal error"}` + "\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.mockBehavior(context.Background(), tt.inputCategory.Category, t)
			handler := handlers.NewProductHandler(mock)

			e := echo.New()
			e.POST("/api/v1/products/", handler.AddCategory)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/products/", bytes.NewBufferString(tt.inputBody))
			req.Header.Set("Content-Type", "application/json")

			e.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedResposeBody, w.Body.String())
			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}

}
