package handlers

import (
	"log"
	"net/http"

	"github.com/Gergenus/commerce/user-service/pkg/jwtpkg"
	"github.com/labstack/echo/v4"
)

type UserMiddleware struct {
	jwtProduct jwtpkg.UserJWTInterface
}

func NewUserMiddleware(jwtProduct jwtpkg.UserJWTInterface) UserMiddleware {
	return UserMiddleware{
		jwtProduct: jwtProduct,
	}
}

func (p UserMiddleware) ConfirmationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("AccessToken")
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "getting token error",
			})
		}
		if cookie.Value == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "No auth token",
			})
		}
		role, uuid, err := p.jwtProduct.ParseToken(cookie.Value)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "Invalid token",
			})
		}
		c.Set("role", role)
		c.Set("uuid", uuid)
		return next(c)
	}
}
