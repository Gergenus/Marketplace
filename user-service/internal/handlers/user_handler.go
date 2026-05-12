package handlers

import (
	"errors"
	"net/http"

	"github.com/Gergenus/commerce/user-service/internal/models"
	"github.com/Gergenus/commerce/user-service/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const AccessTokenDuration = 60 * 24 * 60 * 60
const RefreshTokenDuration = 60 * 24 * 60 * 60

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
		if errors.Is(err, service.ErrIncorrectEmail) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "incorrect email",
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

func (u *UserHandler) Login(c echo.Context) error {
	var loginReq models.LoginRequest

	err := c.Bind(&loginReq)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid payload",
		})
	}
	ip := c.RealIP()
	AccessToken, RefreshToken, err := u.srv.Login(c.Request().Context(), loginReq.Email, loginReq.Password, c.Request().UserAgent(), ip)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}
	err = setCookie(c, "AccessToken", AccessToken, AccessTokenDuration)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	err = setCookie(c, "RefreshToken", RefreshToken, RefreshTokenDuration)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
	return c.JSON(http.StatusAccepted, map[string]interface{}{
		"AccessToken":  AccessToken,
		"RefreshToken": RefreshToken,
	})
}

func setCookie(c echo.Context, key, value string, duration int) error {
	cookie := http.Cookie{
		Name:     key,
		HttpOnly: true,
		Secure:   true,
		MaxAge:   duration,
		Value:    value,
		Path:     "/api/v1",
	}
	c.SetCookie(&cookie)
	return nil
}

func (u *UserHandler) Refresh(c echo.Context) error {
	oldRefresh, err := c.Cookie("RefreshToken")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal server error",
		})
	}
	refreshUUID, err := uuid.Parse(oldRefresh.Value)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal server error",
		})
	}
	expiredAccesToken, err := c.Cookie("AccessToken")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal server error",
		})
	}
	newRefresh, newAccess, err := u.srv.RefreshToken(c.Request().Context(), refreshUUID, c.Request().UserAgent(), c.RealIP(), expiredAccesToken.Value)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRefreshSession) {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "invalid refresh session",
			})
		}
		if errors.Is(err, service.ErrTokenExpired) {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "token expired",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal server error",
		})
	}
	setCookie(c, "RefreshToken", newRefresh.String(), RefreshTokenDuration)
	setCookie(c, "AccessToken", newAccess, AccessTokenDuration)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"AccessToken":  newAccess,
		"RefreshToken": newRefresh,
	})
}

func (u *UserHandler) Logout(c echo.Context) error {
	cookie, err := c.Cookie("RefreshToken")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "bad request",
		})
	}
	err = u.srv.Logout(c.Request().Context(), cookie.Value)
	if err != nil {
		if errors.Is(err, service.ErrNoSessionFound) {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "unauthorized",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "internal error",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"logout": cookie.Value,
	})
}

func (u *UserHandler) Verification(c echo.Context) error {
	verifToken := c.QueryParam("token")
	if verifToken == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "invalid request",
		})
	}

	err := u.srv.Verification(c.Request().Context(), verifToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "unauthorized",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "success",
	})
}
