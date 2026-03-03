package handlers

import (
	"net/http"

	"github.com/Gergenus/commerce/user-service/internal/models"
	"github.com/Gergenus/commerce/user-service/internal/service"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	service service.ServiceInterface
}

func NewProductHandler(service service.ServiceInterface) ProductHandler {
	return ProductHandler{
		service: service,
	}
}

// POST request
func (p *ProductHandler) AddCategory(c echo.Context) error {
	var category models.Category
	err := c.Bind(&category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err,
		})
	}

	id, err := p.service.AddCategory(c.Request().Context(), category.Category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err,
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
