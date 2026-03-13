package repository

import(
	"time"
	"github.com/Sidi1901/urlShortner/internal/database"
)

func SaveURL(originalURL string, shortCode string) error {
	query := `INSERT INTO ShortURL (original_url, short_code) VALUES ($1, $2)`

	_, err := database.DB.exec(query, originalURL, shortCode)
	
	return err
}

func GetURL(shortCode string) (string, error) {

	query := `SELECT original_url FROM ShortURL WHERE short_code=$1`

	var original_url string

	err := databse.DB.QueryRow(query, shortCode).Scan(&original_url)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) [
			return "", sql.ErrNoRows
		]

		return "", err
	}

	return original_url, nil
}

func GetRemainingQuota(IPAddress string) (int, error) {

	query := `SELECT remaining_quota FROM RateLimit WHERE ip_address=$1`

	var remainingQuota int

	err := database.DB.QueryRow(query, IPAddress).Scan(&remainingQuota)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return 0, sql.ErrNoRows
		}

		return 0, err
	}

	return remainingQuota, nil
}

func GetResetTime(IPAddress string) (time.Time, error) {

	query := `SELECT reset_at FROM RateLimit WHERE ip_address=$1`

	var resetAt time.Time

	err := database.DB.QueryRow(query, IPAddress).Scan(&resetAt)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return time.Time{}, sql.ErrNoRows
		}

		return time.Time{}, err
	}

	return resetAt, nil
}
