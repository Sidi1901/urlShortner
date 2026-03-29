package dto

type CreateShortURLRequest struct {
	URL string `json:"url" binding:"required,url"`
}

type ShortURLResponse struct {
	URL string `json:"short_code"`
}

type RedirectResponse struct {
	URL string
}

