package service

import(
   "time"
   "strconv"
   "os"
   "fmt"
   "github.com/Sidi1901/urlShortner/internal/repository"
   repo "github.com/Sidi1901/urlShortner/pkg/utils"
   "github.com/asaskevich/govalidator"
   "github.com/google/uuid"
   "github.com/gin-gonic/gin"
)


resp := response{
	
}



func CreateShortURL(url string, ip string, expiry time.Duration, CustomShortCode string) (resp,error) {
	

	// 1. Implement rate limiting

	// Try to get the value of the IP address from the database. If the value does not exist, it means that it's the first time the IP address is making a request, so we set the rate limit for that IP address. If the value exists, we check if it's less than or equal to 0, which means that the rate limit has been exceeded. If the rate limit is exceeded, we return a 503 Service Unavailable status with a message indicating that the rate limit has been exceeded and when it will reset.

	quota, err := repo.GetRemainingQuota(c.IP())

	if err == nil{
		// Rate limit exists for the IP address. Check if it's exceeding the quota. 
	
		// Check if val <=0 indicating expired key in db which was initially set to 30 minues.
		
		if val <= 0 {
			// limit is of type time.Duration and the value is stored in nanoseconds. We convert it to minutes by dividing it by time.Minute 
			limit, _ := repo.GetResetTime(c.IP())
			return "", fmt.Errorf("Rate limit is exceeded. Try again in %v minutes", limit / time.Minute)
		}

		return "", mt.Errorf("Oops! Some error occurred in short url service")
	}

	// 2. check if the input is an actual URL
	if !govalidator.IsURL(body.URL) {
		return "", fmt.Errorf()
	}

	// 3. Check for domain error
	if !helpers.IsValidDomain(body.URL) {
		return "", fmt.Errorf()
	}

	// 4. enforce ssl for https

	body.URL = helpers.EnforceHTTP(body.URL)



	/*
		After all checks have been passed, Create (or input from user) unique ID for url = domain + customShortCode.
		Check unique url is not already exists in DB as well.

	*/

	id := uuid.New().String()[:6]


	if CustomShortCode == "" {
		customShortCode = uuid.New().String()[:6]
	}

	_, err := repo.GetURL(customShortCode)

	if err == nil {
		return "", fmt.Errorf("Custom short url is already in use")
	}

	var id string

	id := uuid.New().String()[:6]
	
	// Save data in table ShortURL
	_ = repo.SaveURL(id, body.URL, customShortCode, expiry, c.IP())

	// Save data in table RateLimit 1) 

	err = repo.SetRateLimit(c.IP(), quota-1, time.Now().Add(30*time.Minute), time.Now())

	if err != nil {
		return "", fmt.Errorf("Unable to connect to the database")
	}

	// Get reset time for sending in response

	nextResetTime, err2 := repo.GetResetTime(c.IP())

	if err2 != nil {
		return "", fmt.Errorf("Unable to connect to the database")
	}
	
	return [customShortCode, quota-1, nextResetTime], nil

}
