package repository

import(
	"time"
	"errors"
	"context"
	"github.com/Sidi1901/urlShortner/internal/database"
	"github.com/jmoiron/sqlx"
)



type ShortURLRepository interface {
	// Save Short code for an original URL
	SaveShortCode(ctx context.Context, shortURL *model.ShortURL) error

	// Read operation
	GetByShortCode(ctx context.Context, code string) (*model.ShortURL, error)

	// Update Short code for an origial URL
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
	query := `INSERT into short_codes(id, short_code, original_url, created_at, expires_at, ip_address) VALUES(:id, :short_code, :original_url, :created_at, :expires_at, :ip_address)`

	//Ensure ID is set

	if sURL.ID == uuid.Nil {
		sURL.ID = uuid.New()
	}

	_,err := r.db.NamedExecContext(ctx, query, sURL)

	return err
}


func (r *shortURLRepository) GetByShortCode(ctx context.Context, code string) (*model.ShortURL, error) {
	query := `SELECT id, short_code, original_url, created_at, expires_at, ip_address FROM short_url WHERE short_code = $1 LIMIT 1`

	var shortURL model.ShortURL
	_,err := r.db.GetContext(ctx, &shortURL, query, code)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return nil, errors.New("short code not found")
		}
	}

	return &shortURL, nil

}


func (r *shortURLRepository) UpdateShortCode(ctx context.Context, URL *model.ShortURL) error {


}

func (r *shortURLRepository) GetByShortCode(ctx context.Context, code string) error {


}
