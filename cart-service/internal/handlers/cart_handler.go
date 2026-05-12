package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Gergenus/commerce/cart-service/internal/models"
	"github.com/Gergenus/commerce/cart-service/internal/service"
	"github.com/Gergenus/commerce/cart-service/proto"
	"github.com/labstack/echo/v4"
)

type CartHandler struct {
	srv service.CartServiceInterface
	proto.UnimplementedOrderServiceServer
}

func NewCartHandler(srv service.CartServiceInterface) *CartHandler {
	return &CartHandler{srv: srv}
}

func (ch *CartHandler) AddToCart(c echo.Context) error {
	UUID := c.Get("uuid")
	UUIDString, ok := UUID.(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}

	var data models.AddCartRequest
	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid payload",
		})
	}

	productId := strconv.Itoa(data.ProductId)
	err = ch.srv.AddToCart(c.Request().Context(), UUIDString, productId, data.Stock)
	if err != nil {
		if errors.Is(err, service.ErrInsufficientStock) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"warn": "insufficient stock",
			})
		}
		if errors.Is(err, service.ErrAlredyAdded) {
			return c.JSON(http.StatusConflict, map[string]interface{}{
				"warn": "product is already added. Delete it and add again",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}

func (ch *CartHandler) DeleteFromCart(c echo.Context) error {
	UUID := c.Get("uuid")
	UUIDString, ok := UUID.(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}
	var data models.DeleteFromCartRequest
	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid payload",
		})
	}
	productId := strconv.Itoa(data.ProductId)
	err = ch.srv.DeleteFromCart(c.Request().Context(), UUIDString, productId)
	if err != nil {
		if errors.Is(err, service.ErrNoCartProductFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"warn": "no cart product found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal server",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "product " + strconv.Itoa(data.ProductId) + " has been deleted",
	})
}

func (ch *CartHandler) Cart(c echo.Context) error {
	UUID := c.Get("uuid")
	UUIDString, ok := UUID.(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}
	carts, err := ch.srv.Cart(c.Request().Context(), UUIDString)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal server",
		})
	}
	return c.JSON(http.StatusOK, carts)
}
