package repository

import(
	"time"
	"errors"
	"context"
	"github.com/Sidi1901/urlShortner/internal/database"
	"github.com/jmoiron/sqlx"
)

type ShortCodeRepository interface {
	// Save Short code for an original URL
	SaveShortCode(ctx context.Context, shortcode *model.ShortCodes) error

	// Get data for example ShortCode by Original URL
	FindByOriginalURl(ctx context.Context, url model.ShortCodes)

	// Update Short code for an origial URL
	UpdateShortCode()

	// Delete Short code foran original URL
	DeleteShortCode()
}

type shortCodeRepository struct {
	db *sqlx.DB
}

func NewURLRepository(db *sql.db) ShortCodeRepository{
	return &shortCodeRepository{db:db}
}

// Getters and setters for ShortCodeRepository

func (r *shortCodeRepository) SaveShortCode (ctx context.Context, shortcodes *model.ShortCodes) error{
	query := `INSERT into short_codes(id, short_code, original_url, created_at, expires_at, ip_address) VALUES($1,$2,$3,$4,$5,$6)`

	_,err := r.db.ExecContext(
		ctx,
		query,
		shortcodes.ID,
		shortcodes.ShortCode,
		shortcodes.CreatedAt,
		shortcodes.ExpiresAt,
		shortcodes.IPAddress,
		
	)

	return err
}







type RatelimitRepository interface {
	CreateQuota()
	GetQuota()
	UpdateQuota()	
}


























































