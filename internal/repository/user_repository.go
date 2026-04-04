package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sidi1901/urlShortner/internal/logger"
	"github.com/Sidi1901/urlShortner/internal/model"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, email string) (model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, email string) error
}

func (r *Repository) CreateUser(ctx context.Context, user *model.User) error {
	query := `INSERT INTO url_shortner.users (email, name, password, created_at, updated_at, user_role, user_type, user_id) VALUES (:email, :name, :password, :created_at, :updated_at, :user_role, :user_type, :user_id)`

	fmt.Println("Creating user with user:", user)

	_, err := r.db.NamedExecContext(ctx, query, user)

	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to save user")
		return fmt.Errorf("Failed to save user: %w", err)
	}

	return nil
}

func (r *Repository) GetUser(ctx context.Context, email string) (model.User, error) {
	query := `SELECT * FROM url_shortner.users WHERE email = :email`

	var usermodel model.User

	params := map[string]interface{}{
		"email": email,
	}

	rows, err := r.db.NamedQueryContext(ctx, query, params)

	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to Get User aata")
		return usermodel, ErrUserNotFound
	}

	defer rows.Close()

	if !rows.Next() {
		return usermodel, ErrUserNotFound
	}

	if err := rows.StructScan(&usermodel); err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to scan user")
		return usermodel, err
	}

	logger.Log.Info("User retrieved successfully")

	return usermodel, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *model.User) error {
	query := `UPDATE url_shortner.users SET
	email = :original_url,
	name = :expiry_duration,
	password = :password,
	user_type = :users_type,
	user_role = :user_role,
	updated_at = :updated_at
	WHERE email = :email`

	_, err := r.db.NamedExecContext(ctx, query, user)

	if err != nil {
		return fmt.Errorf("Failed to update user - %w", err)
	}

	logger.Log.Info("User updated successfully")

	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, email string) error {
	query := `DELETE FROM url_shortner.users WHERE email = :email`

	params := map[string]interface{}{
		"email": email,
	}

	_, err := r.db.NamedExecContext(ctx, query, params)

	if err != nil {
		return fmt.Errorf("Failed to delete user - %w", err)
	}

	logger.Log.WithFields(map[string]interface{}{
		"Email": email,
	}).Info("User deleted successfully")

	return nil
}
