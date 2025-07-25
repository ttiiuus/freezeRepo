package postgres

import (
	"auth/internal/entity"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

/*
Create(user *User) error
GetByID(id uuid.UUID) (*User, error)
GetByUsername(username string) (*User, error)
Update(user *User) error
Delete(id uuid.UUID) error
List() ([]*User, error)
Exists(username string) (bool, error)
*/

type UserPostgres struct {
	pool *pgxpool.Pool
}

func NewUserPostgres(pool *pgxpool.Pool) *UserPostgres {
	return &UserPostgres{pool: pool}
}

func (pu *UserPostgres) Create(user *entity.User) error {
	query := `
	INSERT INTO users (id, username, email, password_hash, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := pu.pool.Exec(context.Background(),
		query, user.ID, user.Username, user.Email, user.PasswordHash, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("cannot create user %v", err)
	}
	return nil
}

func (pu *UserPostgres) GetByUsername(username string) (*entity.User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE username = $1`

	row := pu.pool.QueryRow(context.Background(), query, username)

	var user entity.User
	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (pu *UserPostgres) GetByID(id uuid.UUID) (*entity.User, error) {
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE id = $1`

	row := pu.pool.QueryRow(context.Background(), query, id)

	var user entity.User
	err := row.Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt,
	)
	//можно добавить обработку ошибок
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (pu *UserPostgres) Update(user *entity.User) error {
	query := `
		UPDATE users
		SET username = $1,
			email = $2,
			password_hash = $3,
			updated_at = NOW()
		WHERE id = $4
	`

	cmdTag, err := pu.pool.Exec(
		context.Background(),
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no user found with id %s", user.ID)
	}

	return nil
}

func (pu *UserPostgres) Delete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	cmdTag, err := pu.pool.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no user found with id %s", id)
	}

	return nil
}

func (pu *UserPostgres) List() ([]*entity.User, error) {
	query := `SELECT username, email FROM users`
	rows, err := pu.pool.Query(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("we cannot list users, pls check database %v", err)
	}
	defer rows.Close()

	var users []*entity.User

	for rows.Next() {
		var user entity.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return users, nil
}

func (pu *UserPostgres) Exists(username string) (bool, error) {
	query := `
	SELECT EXISTS (SELECT username FROM users WHERE username = $1)
	`
	var exists bool
	err := pu.pool.QueryRow(context.Background(), query, username).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}

	return exists, nil
}

func (pu *UserPostgres) UpdateBalance(userID uuid.UUID, balance float64) error {
	query := `
	UPDATE users
	SET balance = $1
	WHERE id = $2
	`
	cmdTag, err := pu.pool.Exec(context.Background(), query, balance, userID)
	if err != nil {
		return fmt.Errorf("failed to update user balance: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no user found with id %s", userID)
	}

	return nil
}
