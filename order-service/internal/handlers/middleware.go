package handlers

import (
	"net/http"

	"github.com/Gergenus/commerce/order-service/pkg/jwt"
	"github.com/labstack/echo/v4"
)

func OrderAuth(next echo.HandlerFunc) echo.HandlerFunc {
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
		role, uuid, err := jwt.ParseToken(cookie.Value)
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
		c.Set("role", role)
		c.Set("uuid", uuid)
		return next(c)
	}
}

func SellerAuth(next echo.HandlerFunc) echo.HandlerFunc {
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
		role, uuid, err := jwt.ParseToken(cookie.Value)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "Invalid token",
			})
		}
		if role != "seller" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "invalid role",
			})
		}
		c.Set("role", role)
		c.Set("uuid", uuid)
		return next(c)
	}
}
