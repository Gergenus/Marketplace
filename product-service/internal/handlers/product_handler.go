package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Gergenus/commerce/product-service/internal/models"
	"github.com/Gergenus/commerce/product-service/internal/service"
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
			"error": "Invalid request payload",
		})
	}

	id, err := p.service.AddCategory(c.Request().Context(), category.Category)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Internal error",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (p *ProductHandler) CreateProduct(c echo.Context) error {
	if c.Get("role") == "" || c.Get("role") != "seller" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "unsuitable role",
		})
	}

	var product models.Product
	err := c.Bind(&product)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request payload",
		})
	}
	sellerID, ok := c.Get("seller_id").(int)
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "claims error",
		})
	}
	product.SellerID = sellerID
	id, err := p.service.CreateProduct(c.Request().Context(), product)
	if err != nil {
		if errors.Is(err, service.ErrMoreThanOneProductInstance) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": service.ErrMoreThanOneProductInstance.Error(),
			})
		}
		if errors.Is(err, service.ErrNoSuchCategoryExists) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error":       service.ErrNoSuchCategoryExists.Error(),
				"category_id": product.CategoryID,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal error",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (p *ProductHandler) GetStockByID(c echo.Context) error {
	productIdString := c.QueryParam("product_id")
	sellerId := c.Get("seller_id").(int)

	productId, err := strconv.Atoi(productIdString)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request payload",
		})
	}
	stock, err := p.service.GetStockByID(c.Request().Context(), productId, sellerId)
	if err != nil {
		if errors.Is(err, service.ErrStockNotFound) {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Stock not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal error",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"product_id": productId,
		"seller_id":  sellerId,
		"stock":      stock,
	})
}

func (p *ProductHandler) AddStockByID(c echo.Context) error {
	var stockReq models.AddStockRequest
	sellerID := c.Get("seller_id").(int)
	err := c.Bind(&stockReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request payload",
		})
	}

	id, err := p.service.AddStockByID(c.Request().Context(), sellerID, stockReq.ProductID, stockReq.Number)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal error",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (p *ProductHandler) GetProductByID(c echo.Context) error {
	id := c.QueryParam("product_id")

	product_id, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request payload",
		})
	}
	product, err := p.service.GetProductByID(c.Request().Context(), product_id)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "Product not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal error",
		})
	}
	return c.JSON(http.StatusOK, product)
}
