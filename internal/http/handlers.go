package http

import (
	"auth/internal/entity"
	"auth/internal/service"
	"auth/internal/usecase"
	"auth/pkg/logger"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	userUsecase   usecase.UserUsecase
	reportUsecase usecase.ReportUsecase
	jwtService    service.JWTService
}

func NewHandler(userUsecase usecase.UserUsecase, reportUsecase usecase.ReportUsecase, jwtService service.JWTService) *Handler {
	return &Handler{
		userUsecase:   userUsecase,
		reportUsecase: reportUsecase,
		jwtService:    jwtService,
	}
}

type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Register(c echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	token, err := h.userUsecase.RegisterUser(req.Username, req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	setAuthCookie(c, token)

	return c.JSON(http.StatusCreated, map[string]string{"message": "registered successfully"})
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func (h *Handler) Login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	token, err := h.userUsecase.LoginUser(req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	setAuthCookie(c, token)

	return c.JSON(http.StatusCreated, map[string]string{"message": "logged in successfully"})
}

func (h *Handler) ListUsers(c echo.Context) error {
	users, err := h.userUsecase.ListUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to list users"})
	}
	return c.JSON(http.StatusOK, users)
}

func (h *Handler) CreateReport(c echo.Context) error {
	var report entity.Report
	if err := c.Bind(&report); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := h.reportUsecase.CreateReport(&report); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create report"})
	}
	return c.JSON(http.StatusCreated, report)
}

func setAuthCookie(c echo.Context, token string) {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   3600,
	}
	c.SetCookie(cookie)
	mess := fmt.Sprintf("sucssesfull set cookie %s", cookie)
	logger.Logger.Info().Msg(mess)

}

func (h *Handler) CheckAuth(c echo.Context) error {
	username, ok := c.Get("username").(string)
	if !ok || username == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok", "username": username})
}

func (h *Handler) GetUserReports(c echo.Context) error {
	userID := c.Param("id")
	uuid_user, err := uuid.Parse(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid user ID"})
	}
	reports, err := h.reportUsecase.GetUserReports(uuid_user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get user reports"})
	}
	return c.JSON(http.StatusOK, reports)
}

func (h *Handler) PurchaseReport(c echo.Context) error {
	reportID := c.Param("report_id")
	if reportID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "report ID is required"})
	}

	// Call the use case to purchase the report
	if err := h.reportUsecase.PurchaseReport(reportID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to purchase report"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "report purchased successfully"})
}

/*
type CookieStruct struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (h *Handler) SetCookie(c echo.Context) error {
	var cookie CookieStruct
	if err := c.Bind(&cookie); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}
	cookieValue := cookie.Value
	if cookie.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cookie name cannot be empty"})
	}
	if cookieValue == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cookie value cannot be empty"})
	}
	http.SetCookie(c.Response().Writer, &http.Cookie{
		Name:  cookie.Name,
		Value: cookie.Value,
	})
	return nil
}
func (h *Handler) GetCookie(c echo.Context) (*http.Cookie, error) {
	var name string
	if err := c.Bind(&name); err != nil {
		return nil, c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}
	if name == "" {
		return nil, c.JSON(http.StatusBadRequest, map[string]string{"error": "Cookie name cannot be empty"})
	}
	cookie, err := c.Cookie(name)
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, fmt.Errorf("cookie %s not found", name)
		}
		return nil, fmt.Errorf("error retrieving cookie %s: %v", name, err)
	}
	return cookie, nil
}

func (h *Handler) DeleteCookie(c echo.Context) error {
	var name string
	if err := c.Bind(&name); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}
	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Cookie name cannot be empty"})
	}
	cookie := &http.Cookie{
		Name:   name,
		Value:  "",
		MaxAge: -1, // Устанавливаем MaxAge в -1 для удаления cookie
	}
	http.SetCookie(c.Response().Writer, cookie)
	return nil
}
*/
