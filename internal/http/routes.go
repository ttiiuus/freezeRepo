package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterRoutes(e *echo.Echo, h *Handler) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://192.169.208.1:8085"}, // upload-сервер
		AllowMethods:     []string{echo.GET, echo.POST},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	e.File("/", "web/auth.html")
	e.Static("/", "web")
	e.POST("/login", h.Login)
	e.POST("/register", h.Register)
	e.POST("/reports", h.CreateReport) //mongodb

	api := e.Group("/api", AuthMiddleware(h.jwtService))

	api.GET("/check", h.CheckAuth)
	api.GET("/users", h.ListUsers)

	api.GET("/:id/reports", h.GetUserReports)                      //mongodb
	api.POST("/api/reports/:report_id/purchase", h.PurchaseReport) //mongodb

	//api.GET("")
	//e.GET("/users", ListUsers)
	//e.GET("/users/:id", GetUser)
	//e.PUT("/users/:id", UpdateUser)
	//e.DELETE("/users/:id", DeleteUser)
}
