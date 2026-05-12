package handlers

import (
	"net"
	"net/http"
	"sync"

	"github.com/Gergenus/commerce/product-service/pkg/jwtpkg"
	ratelimiter "github.com/Gergenus/rateLimiter"
	"github.com/labstack/echo/v4"
)

type ProductMiddleware struct {
	jwtProduct       jwtpkg.JWTinterface
	RateLimitingUser map[string]ratelimiter.RateLimiter
}

func NewProductMiddleware(jwtProduct jwtpkg.JWTinterface) ProductMiddleware {
	return ProductMiddleware{
		jwtProduct:       jwtProduct,
		RateLimitingUser: make(map[string]ratelimiter.RateLimiter),
	}
}

func (p ProductMiddleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
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

func (p ProductMiddleware) CreateProductAuth(next echo.HandlerFunc) echo.HandlerFunc {
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
		role, uuid, err := p.jwtProduct.ParseToken(cookie.Value)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "Invalid token",
			})
		}

		if role == "" || role != "seller" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "unsuitable role",
			})
		}

		c.Set("role", role)
		c.Set("uuid", uuid)
		return next(c)
	}
}

func (p ProductMiddleware) RateLimiting(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip, err := getIP(c.Request().RemoteAddr)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "internal error",
			})
		}
		var mu sync.Mutex
		mu.Lock()
		if userLimiter, ok := p.RateLimitingUser[ip]; ok {
			if userLimiter.Allow() {
				return next(c)
			} else {
				return c.JSON(http.StatusTooManyRequests, map[string]string{
					"error": "too many requests",
				})
			}
		} else {
			p.RateLimitingUser[ip] = ratelimiter.NewRateLimiter(ratelimiter.NewLimit(50), 20)
		}
		mu.Unlock()
		return next(c)
	}
}

func getIP(addr string) (string, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}
	return host, nil
}
