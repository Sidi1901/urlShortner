package repository

/*
Queries for
user has which shorts urls for which long urls and check if that is active
user long url short url is_active expiry
*/

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UserURL struct {
	ID        int
	UserID    int
	ShortCode string
	LongURL   string
	CreatedAt string
}

type UserURLRepository interface {
	GetURLByUserID(ctx context.Context, user_id int) ([]UserURL, error)
}

type userURLRepository struct {
	db *sqlx.DB
}

func NewUserURLRepository(db *sqlx.DB) UserURLRepository {
	return &userURLRepository{db: db}
}

func (r *userURLRepository) GetURLByUserID(ctx context.Context, userID int) ([]UserURL, error) {
	query := `
		SELECT id, user_id, short_code, long_url, created_at
		FROM short_urls
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urls []UserURL

	for rows.Next() {
		var u UserURL
		if err := rows.Scan(&u.ID, &u.UserID, &u.ShortCode, &u.LongURL, &u.CreatedAt); err != nil {
			return nil, err
		}
		urls = append(urls, u)
	}

	return urls, nil
}
