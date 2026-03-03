package handlers

import (
	"errors"
	"net/http"

	"github.com/Gergenus/commerce/user-service/internal/models"
	"github.com/Gergenus/commerce/user-service/internal/service"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	srv service.UserServiceInterface
}

func NewUserHandler(srv service.UserServiceInterface) UserHandler {
	return UserHandler{srv: srv}
}

func (u *UserHandler) Register(c echo.Context) error {
	var user models.User

	err := c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid payload",
		})
	}
	uid, err := u.srv.AddUser(c.Request().Context(), user)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "user already exists",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"uid": uid.String(),
	})
}
