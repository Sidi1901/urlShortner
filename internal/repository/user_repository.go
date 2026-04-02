package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sidi1901/urlShortner/internal/logger"
	"github.com/Sidi1901/urlShortner/internal/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, email string) (model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, email string) error
}

func (r *Repository) CreateUser(ctx context.Context, user *model.User) error {
	query := `INSERT INTO url_shortner.user (id, email, name, password, token, refresh_token, created_at, updated_at, user_role, user_type, user_id) VALUES (:id, :email, :name, :password, :token, :refresh_token, :created_at, :updated_at, :user_role, :user_type, :user_id)`

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
	query := `SELECT * FROM url_shortner.user WHERE email = :email`

	var usermodel model.User

	params := map[string]interface{}{
		"email": email,
	}

	rows, err := r.db.NamedQueryContext(ctx, query, params)

	if err != nil {
		logger.Log.WithFields(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to Get User aata")
		return usermodel, fmt.Errorf("Failed to Get User data %w", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return usermodel, errors.New("User not found")
	}

	logger.Log.WithFields(map[string]interface{}{
		"ID": usermodel.ID,
	}).Info("User retrieved successfully")

	return usermodel, nil
}

func (r *Repository) UpdateUser(ctx context.Context, user *model.User) error {
	query := `UPDATE`

	_, err := r.db.NamedExecContext(ctx, query, user)

	if err != nil {
		return fmt.Errorf("Failed to update user - %w", err)
	}

	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, email string) error {
	query := `DELETE`

	params := map[string]interface{}{
		"email": email,
	}

	_, err := r.db.NamedExecContext(ctx, query, params)

	if err != nil {
		return fmt.Errorf("Failed to delete user - %w", err)
	}

	return nil
}
