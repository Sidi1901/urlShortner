package handler

import(
	"fmt"
	repo "github.com/Sidi1901/urlShortner/internal/repository"
)

func GetOriginalURL (code string) (string, error){

	val, err := repo.GetURL(code)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("URL not found")
		}
		return "", fmt.Errorf("Unable to connect to the database")
	}

	return val, nil
}