package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (u *UserHandler) RegistrationConfirmation(c echo.Context) error {
	uuid, ok := c.Get("uuid").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "unauthorized",
		})
	}
	err := u.srv.ConfirmUser(c.Request().Context(), uuid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}
