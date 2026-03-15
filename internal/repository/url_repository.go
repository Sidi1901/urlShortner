package repository

import(
	"time"
	"github.com/Sidi1901/urlShortner/internal/database"
)

// Set into ShortURL
func SaveURL(id string, originalURL string, shortCode string, expiry time.Duration, ipAddress string) error {
	query := `INSERT INTO ShortURL (id, short_code, original_url, created_at, expires_at, created_by) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := database.DB.exec(query, id, shortCode, originalURL, time.Now(), time.Now().Add(expiry), "system")
	
	return err
}

// Get from ShortURL
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


// Set into RateLimit
func SetRateLimit(IPAddress string, remainingQuota int, resetAt time.Time) error {

	query := `INSERT INTO RateLimit (ip_address, remaining_quota, reset_at, updated_at) VALUES ($1, $2, $3, $4) ON CONFLICT (ip_address) DO UPDATE SET remaining_quota=$2, reset_at=$3, updated_at=$4`	

	updatedAt := time.Now()
	_, err := database.DB.Exec(query, IPAddress, remainingQuota, resetAt, updatedAt)

	return err
}


// Get from RateLimit
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


//Get from RateLimit
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


