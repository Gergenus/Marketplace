package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Gergenus/commerce/product-service/internal/models"
	"github.com/Gergenus/commerce/product-service/internal/service"
	"github.com/Gergenus/commerce/product-service/proto"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	service service.ServiceInterface
	proto.UnimplementedAvailablilityServiceServer
	proto.UnimplementedOrderServiceServer
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
		if errors.Is(err, service.ErrCategoryAlreadyExists) {
			return c.JSON(http.StatusBadRequest, map[string]any{
				"error": "category already exists",
			})
		}
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Internal error",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

func (p *ProductHandler) CreateProduct(c echo.Context) error {
	var product models.Product
	err := c.Bind(&product)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request payload",
		})
	}
	sellerID, ok := c.Get("uuid").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "unauthorized",
		})
	}
	product.SellerID = sellerID
	id, err := p.service.CreateProduct(c.Request().Context(), product)
	if err != nil {
		if errors.Is(err, service.ErrMoreThanOneProductInstance) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "more than one instance of a product",
			})
		}
		if errors.Is(err, service.ErrNoSuchCategoryExists) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "no such category exists",
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
	sellerId := c.Get("uuid").(string)

	productId, err := strconv.Atoi(productIdString)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request payload",
		})
	}
	stock, err := p.service.GetStockByID(c.Request().Context(), productId)
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
	sellerID := c.Get("uuid").(string)
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

func (p *ProductHandler) Products(c echo.Context) error {
	query := c.QueryParam("search_query")
	if query == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "bad request",
		})
	}
	page := c.QueryParam("page")
	if page == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "bad request",
		})
	}
	pageSize := c.QueryParam("page_size")
	if pageSize == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "bad request",
		})
	}
	products, err := p.service.Products(c.Request().Context(), query, page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal error",
		})
	}
	return c.JSON(http.StatusOK, products)
}
