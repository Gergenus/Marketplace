package tests

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gergenus/commerce/product-service/internal/handlers"
	"github.com/Gergenus/commerce/product-service/internal/models"
	"github.com/Gergenus/commerce/product-service/internal/service"
	"github.com/Gergenus/commerce/product-service/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct(t *testing.T) {
	tests := []struct {
		name                string
		inputBody           string
		inputProduct        models.Product
		userUUID            uuid.UUID
		mockBehavior        func(ctx context.Context, product models.Product, id uuid.UUID, t *testing.T) *mocks.MockServiceInterface
		expectedCode        int
		expectedResonseBody string
	}{
		{
			name:         "OK",
			inputBody:    `{"product_name":"table", "price": 300.0, "category_id": 1}`,
			userUUID:     uuid.New(),
			inputProduct: models.Product{ProductName: "table", Price: 300.0, CategoryID: 1},
			mockBehavior: func(ctx context.Context, product models.Product, id uuid.UUID, t *testing.T) *mocks.MockServiceInterface {
				mock := mocks.NewMockServiceInterface(t)
				product.SellerID = id.String()
				mock.EXPECT().CreateProduct(ctx, product).Return(1, nil)
				return mock
			},
			expectedCode:        http.StatusOK,
			expectedResonseBody: `{"id":1}` + "\n",
		},
		{
			name:         "unauthorized",
			inputBody:    `{"product_name":"table", "price": 300.0, "category_id": 1}`,
			userUUID:     uuid.Nil,
			inputProduct: models.Product{ProductName: "table", Price: 300.0, CategoryID: 1},
			mockBehavior: func(ctx context.Context, product models.Product, id uuid.UUID, t *testing.T) *mocks.MockServiceInterface {
				mock := mocks.NewMockServiceInterface(t)
				return mock
			},
			expectedCode:        http.StatusUnauthorized,
			expectedResonseBody: `{"error":"unauthorized"}` + "\n",
		},
		{
			name:         "invalid json",
			inputBody:    `{"product_name""table", "price": 300.0, "category_id": 1}`,
			userUUID:     uuid.New(),
			inputProduct: models.Product{ProductName: "table", Price: 300.0, CategoryID: 1},
			mockBehavior: func(ctx context.Context, product models.Product, id uuid.UUID, t *testing.T) *mocks.MockServiceInterface {
				mock := mocks.NewMockServiceInterface(t)
				return mock
			},
			expectedCode:        http.StatusBadRequest,
			expectedResonseBody: `{"error":"Invalid request payload"}` + "\n",
		},
		{
			name:         "many instances",
			inputBody:    `{"product_name": "table", "price": 300.0, "category_id": 1}`,
			userUUID:     uuid.New(),
			inputProduct: models.Product{ProductName: "table", Price: 300.0, CategoryID: 1},
			mockBehavior: func(ctx context.Context, product models.Product, id uuid.UUID, t *testing.T) *mocks.MockServiceInterface {
				mock := mocks.NewMockServiceInterface(t)
				product.SellerID = id.String()
				mock.EXPECT().CreateProduct(ctx, product).Return(-1, service.ErrMoreThanOneProductInstance)
				return mock
			},
			expectedCode:        http.StatusBadRequest,
			expectedResonseBody: `{"error":"more than one instance of a product"}` + "\n",
		},
		{
			name:         "no category exists",
			inputBody:    `{"product_name": "table", "price": 300.0, "category_id": 2}`,
			userUUID:     uuid.New(),
			inputProduct: models.Product{ProductName: "table", Price: 300.0, CategoryID: 2},
			mockBehavior: func(ctx context.Context, product models.Product, id uuid.UUID, t *testing.T) *mocks.MockServiceInterface {
				mock := mocks.NewMockServiceInterface(t)
				product.SellerID = id.String()
				mock.EXPECT().CreateProduct(ctx, product).Return(-1, service.ErrNoSuchCategoryExists)
				return mock
			},
			expectedCode:        http.StatusBadRequest,
			expectedResonseBody: `{"error":"no such category exists"}` + "\n",
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			mock := testCase.mockBehavior(context.Background(), testCase.inputProduct, testCase.userUUID, t)
			handler := handlers.NewProductHandler(mock)

			e := echo.New()

			e.POST("/api/v1/products/create", handler.CreateProduct, func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					if testCase.userUUID != uuid.Nil {
						c.Set("uuid", testCase.userUUID.String())
					}
					return next(c)
				}
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/v1/products/create", bytes.NewBufferString(testCase.inputBody))
			req.Header.Set("Content-Type", "application/json")

			e.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedCode, w.Code)
			assert.Equal(t, testCase.expectedResonseBody, w.Body.String())

		})
	}
}
