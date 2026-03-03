package handlers

import (
	"net/http"
	"strings"

	"github.com/Gergenus/commerce/product-service/pkg/jwtpkg"
	"github.com/labstack/echo/v4"
)

type ProductMiddleware struct {
	jwtProduct jwtpkg.JWTinterface
}

func NewProductMiddleware(jwtProduct jwtpkg.JWTinterface) ProductMiddleware {
	return ProductMiddleware{
		jwtProduct: jwtProduct,
	}
}

func (p ProductMiddleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "No auth token",
			})
		}
		token := strings.Split(tokenString, " ")[1]
		role, sellerID, err := p.jwtProduct.ParseToken(token)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "Invalid token",
			})
		}
		c.Set("role", role)
		c.Set("seller_id", sellerID)
		return next(c)
	}
}
