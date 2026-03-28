package repository

import(
	"fmt"
	"time"
	"errors"
	"context"
	"github.com/jmoiron/sqlx"
)


type ShortURLRepository interface {
	// Save Short code for an original URL
	SaveShortCode(ctx context.Context, shortURL *model.ShortURL) error

	// Read operation
	GetByShortCode(ctx context.Context, code string) (*model.ShortURL, error)

	// Update expiry time for a short code
	UpdateShortCode(ctx context.Context, shortURL *model.ShortURL) error

	// Delete Short code foran original URL
	DeleteShortCode(ctx context.Context, shortCode string) error
}

type shortURLRepository struct {
	db *sqlx.DB
}

func NewURLRepository(db *sqlx.db) ShortURLRepository{
	return &shortURLRepository{db:db}
}

// Getters and setters for ShortURLRepository

func (r *shortURLRepository) SaveShortCode (ctx context.Context, sURL *model.ShortURL) error{
	query := `INSERT into short_url (short_code, original_url, created_at, expires_at, ip_address) VALUES($1, $2, $3, $4, $5)`


	_,err := r.db.GetContext(ctx, query, sURL.ShortCode, sURL.OriginalURL, sURL.CreatedAt, sURL.ExpiresAt, sURL.IPAddress)

	fmt.Printf("Short code saved successfully => %s\n", sURL)

	return err
}


func (r *shortURLRepository) GetByShortCode(ctx context.Context, code string) (*model.ShortURL, error) {
	query := `SELECT short_code, original_url, created_at, expires_at, ip_address FROM short_url WHERE short_code = $1`

	var shortURL model.ShortURL
	_,err := r.db.GetContext(ctx, &shortURL, query, code)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return nil, errors.New("short code not found")
		}
	}

	fmt.Printf("Short code retrieved successfully => %s\n", shortURL)

	return &shortURL, nil

}


func (r *shortURLRepository) UpdateShortCode(ctx context.Context, sURL *model.ShortURL) error {
	query:=`UPDATE short_url SET expires_at = $1 WHERE short_code = $2`
	
	_,err := r.db.ExecContext(ctx, query, sURL.ExpiresAt, sURL.ShortCode)
	if err != nil {
		return err
	}

	return nil
}

func (r *shortURLRepository) DeleteShortCode(ctx context.Context, code string) error {
	query := `DELETE FROM short_url WHERE short_code = $1`

	_, err := r.db.ExecContext(ctx, query, code)

	return err
}
