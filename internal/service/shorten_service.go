package service

import(
"time"
"fmt"
repo "github.com/Sidi1901/urlShortner/internal/repository"
"github.com/Sidi1901/urlShortner/pkg/utils"
"github.com/asaskevich/govalidator"
"github.com/google/uuid"
)


type ShortURLData struct {
	CustomShortCodeURL string  		    `json:"custom_short_code_url,omitempty"`
	Quota		   	   int				`json:"quota,omitempty"`
	ResetTime          time.Duration		`json:"next_reset_time,omitempty"`
}



func CreateShortURL(url string, ip string, expiry time.Duration, CustomShortCode string) (ShortURLData,error) {
	


	// 1. Check if it is an acual URL
	if !govalidator.IsURL(url) {
		return ShortURLData{}, fmt.Errorf("Invalid URL")
	}

	// 2. Check for domain error
	if !utils.IsValidDomain(url) {
		return ShortURLData{}, fmt.Errorf("Invalid Domain")
	}

	// 3. enforce ssl for https

	url = utils.EnforceHTTP(url)



	/*
		After all checks have been passed, Create (or input from user) unique Custom short code for url = domain + customShortCode.
		Check unique Custom short code is not already exists in DB as well.

	*/


	if CustomShortCode == "" {
		CustomShortCode = uuid.New().String()[:6]
	}

	if _,err = repo.GetByShortCode(CustomShortCode); err != nil {
		return ShortURLData{}, fmt.Errorf("Custom short url is already in use")
	}


	
	// Save data in table ShortURL
	_ = repo.SaveShortCode(id, url, CustomShortCode, expiry, ip)

	// Save data in table RateLimit 1) 

	if err = repo.SetRateLimit(ip, quota-1, time.Now().Add(30*time.Minute)); err != nil {
		return ShortURLData{}, fmt.Errorf("Unable to connect to the database")
	}


	// Get reset time for sending in response

	ResetTime, err2 := repo.GetResetTime(ip)
	

	if err2 != nil {
		return ShortURLData{}, fmt.Errorf("Unable to connect to the database")
	}

	RemainingResetTime := time.Until(ResetTime)

	Result := ShortURLData{
		CustomShortCodeURL : CustomShortCode,
		Quota : quota-1,
		ResetTime : RemainingResetTime,
	}
	
	return Result, nil

}
