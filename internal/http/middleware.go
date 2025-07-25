package http

import (
	"auth/internal/service"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(jwtService service.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing or invalid Authorization header"})
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwtService.ValidateJWT(tokenStr)
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token claims"})
			}

			// можно положить username в контекст
			c.Set("uusername", claims["username"])
			return next(c)
		}
	}
}

func AuthMiddleware(jwtService service.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("token")
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing auth token"})
			}
			tokenStr := cookie.Value

			token, err := jwtService.ValidateJWT(tokenStr)
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token claims"})
			}

			username, ok := claims["username"].(string)
			if !ok || username == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "username not found in token"})
			}

			c.Set("username", username) //  username в контексте
			return next(c)
		}
	}
}
