package http

import (
	"auth/config"
	"auth/internal/service"
	"auth/internal/usecase"
	"fmt"

	"github.com/labstack/echo/v4"
)

func StartServer(cfg *config.Config, userUC usecase.UserUsecase, reportUC usecase.ReportUsecase) error {
	e := echo.New()
	jwtService := service.NewJWTService(cfg.Auth.JWTSecret)

	handler := NewHandler(userUC, reportUC, jwtService)
	RegisterRoutes(e, handler)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	fmt.Println("Starting server on", addr)
	return e.Start(addr)
}
