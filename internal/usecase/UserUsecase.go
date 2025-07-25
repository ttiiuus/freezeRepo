package usecase

import (
	"auth/internal/entity"
	"auth/internal/service"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(user *entity.User) error
	RegisterUser(username, email, password string) (string, error)
	LoginUser(username, password string) (string, error)
	GetByUsername(username string) (*entity.User, error)
	Exists(username string) (bool, error)
	ListUsers() ([]*entity.User, error)
	UpdateBalance(userID uuid.UUID, amount float64) error
}

type ReportUsecase interface {
	CreateReport(report *entity.Report) error
	GetUserReports(userID uuid.UUID) ([]*entity.Report, error)
	SetAnonimousIdReport(clientGeneratedID string, userID uuid.UUID) error
	GetUserIdAndPriceByReportId(reportID string) (uuid.UUID, float64, error)
	PurchaseReport(reportID string) error
}

type UserUsecase interface {
	RegisterUser(username, email, password string) (string, error)
	LoginUser(username, password string) (string, error)
	ListUsers() ([]*entity.User, error)
}

type reportUsecase struct {
	reportRepo entity.ReportRepository
	userRepo   entity.UserRepository
}

func NewReportUsecase(reportRepo entity.ReportRepository, userRepo entity.UserRepository) *reportUsecase {
	return &reportUsecase{
		reportRepo: reportRepo,
		userRepo:   userRepo,
	}
}

type userUsecase struct {
	userRepo   entity.UserRepository   // постгрес
	reportRepo entity.ReportRepository // монга
	jwtService service.JWTService
}

func NewUserUsecase(userRepo entity.UserRepository, reportRepo entity.ReportRepository, jwtService service.JWTService) *userUsecase {
	return &userUsecase{
		userRepo:   userRepo,
		reportRepo: reportRepo,
		jwtService: jwtService,
	}
}

func (u *userUsecase) RegisterUser(username, email, password string) (string, error) {
	exists, err := u.userRepo.Exists(username)
	if err != nil {
		return "", err
	}
	if exists {
		return "", fmt.Errorf("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &entity.User{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	//здесь мы создаем пользователя в бд
	err = u.userRepo.Create(user)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %v", err)
	}
	//здесь мы должны создать jwt токен для пользователя

	token, err := u.jwtService.CreateJWT(user)

	if err != nil {
		return "", fmt.Errorf("failed to create JWT token: %v", err)
	}

	return token, nil
}

func (u *userUsecase) LoginUser(username, password string) (string, error) {
	user, err := u.userRepo.GetByUsername(username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	return u.jwtService.CreateJWT(user)
}

func (u *userUsecase) ListUsers() ([]*entity.User, error) {
	users, err := u.userRepo.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %v", err)
	}
	return users, nil
}

//func (u *userUsecase) IssueJWT(username string) (string, error) {
//return service.CreateJWT(username)
//}

func (r *reportUsecase) CreateReport(report *entity.Report) error {
	if err := r.reportRepo.CreateReport(report); err != nil {
		return fmt.Errorf("failed to create report: %v", err)
	}
	return nil
}

func (r *reportUsecase) GetUserReports(userID uuid.UUID) ([]*entity.Report, error) {
	reports, err := r.reportRepo.GetUserReports(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user reports: %v", err)
	}
	return reports, nil
}

func (r *reportUsecase) SetAnonimousIdReport(clientGeneratedID string, userID uuid.UUID) error {
	if err := r.reportRepo.SetAnonimousIdReport(clientGeneratedID, userID); err != nil {
		return fmt.Errorf("failed to set anonymous ID for report: %v", err)
	}
	return nil
}

func (r *reportUsecase) PurchaseReport(reportID string) error {
	userID, price, err := r.reportRepo.GetUserIdAndPriceByReportId(reportID)
	if err != nil {
		return fmt.Errorf("step 1 failed: %w", err)
	}

	if userID == uuid.Nil {
		return fmt.Errorf("no user found for report ID %s", reportID)
	}

	if err := r.userRepo.UpdateBalance(userID, price); err != nil {
		return fmt.Errorf("step 2 failed: %w", err)
	}

	if err := r.reportRepo.PurchaseReport(reportID); err != nil {
		// откат баланса через userRepo
		if rollbackErr := r.userRepo.UpdateBalance(userID, -price); rollbackErr != nil {
			return fmt.Errorf("step 3 failed: %v (also failed to rollback balance: %v)", err, rollbackErr)
		}
		return fmt.Errorf("step 3 failed, rolled back balance: %w", err)
	}

	return nil
}

func (r *reportUsecase) GetUserIdAndPriceByReportId(reportID string) (uuid.UUID, float64, error) {
	return r.reportRepo.GetUserIdAndPriceByReportId(reportID)
}
