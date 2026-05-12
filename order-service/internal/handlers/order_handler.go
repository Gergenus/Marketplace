package handlers

import (
	"errors"
	"net/http"

	"github.com/Gergenus/commerce/order-service/dto"
	"github.com/Gergenus/commerce/order-service/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	srv service.OrderServiceInterface
}

func NewOrderHandler(srv service.OrderServiceInterface) OrderHandler {
	return OrderHandler{srv: srv}
}

func (o OrderHandler) CreateOrder(c echo.Context) error {
	uid := c.Get("uuid").(string)

	var data dto.DeliveryAddress

	err := c.Bind(&data)
	if err != nil || data.Address == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid payload",
		})
	}

	orderID, err := o.srv.CreateOrder(c.Request().Context(), uuid.MustParse(uid), data.Address)
	if err != nil {
		if errors.Is(err, service.ErrOrderNotReserved) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "insufficient stock",
			})
		}
		if errors.Is(err, service.ErrNoCartFound) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "no cart found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"order_id": orderID,
	})
}

func (o OrderHandler) Orders(c echo.Context) error {
	uidString := c.Get("uuid").(string)
	uid := uuid.MustParse(uidString)
	products, err := o.srv.Orders(c.Request().Context(), uid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal error",
		})
	}
	return c.JSON(http.StatusOK, products)
}
