package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	GetByUsername(username string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
	List() ([]*User, error)
	Exists(username string) (bool, error)
	UpdateBalance(userID uuid.UUID, balance float64) error
}

type ReportRepository interface {
	CreateReport(report *Report) error
	GetUserReports(userID uuid.UUID) ([]*Report, error)
	SetAnonimousIdReport(clientGeneratedID string, userID uuid.UUID) error
	GetUserIdAndPriceByReportId(reportID string) (uuid.UUID, float64, error)
	PurchaseReport(reportID string) error
}

type JWT struct {
	Username  string
	JWTSecret string
}

type Report struct {
	Client_generated_id string    `json:"client_generated_id"`
	User_id             string    `json:"user_id"`
	Description         string    `json:"description"`
	Created_at          time.Time `json:"created_at"`
	Report_id           string    `json:"report_id"`
	Is_purchased        bool      `json:"is_purchased"`
	Price               float64   `json:"price"`
}
