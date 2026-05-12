package handlers

import (
	"net/http"

	"github.com/Gergenus/commerce/cart-service/pkg/jwtpkg"
	"github.com/labstack/echo/v4"
)

type CartMiddleware struct {
	jwtProduct jwtpkg.CartJWTInterface
}

func NewCartMiddleware(jwtProduct jwtpkg.CartJWTInterface) CartMiddleware {
	return CartMiddleware{
		jwtProduct: jwtProduct,
	}
}

func (p CartMiddleware) CartMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("AccessToken")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "getting token error",
			})
		}
		if cookie.Value == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "No auth token",
			})
		}

		role, uuid, ver, err := p.jwtProduct.ParseToken(cookie.Value)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "Invalid token",
			})
		}

		if role != "customer" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "invalid role",
			})
		}

		if !ver {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "unverified user",
			})
		}

		c.Set("uuid", uuid)
		return next(c)
	}
}
