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
		

		// 1. Implement rate limiting

		// Try to get the value of the IP address from the database. If the value does not exist, it means that it's the first time the IP address is making a request, so we set the rate limit for that IP address. If the value exists, we check if it's less than or equal to 0, which means that the rate limit has been exceeded. If the rate limit is exceeded, we return a 503 Service Unavailable status with a message indicating that the rate limit has been exceeded and when it will reset.

		quota, err := repo.GetRemainingQuota(ip)

		if err == nil{
			// Rate limit exists for the IP address. Check if it's exceeding the quota. 
		
			// Check if val <=0 indicating expired key in db which was initially set to 30 minues.
			
			if quota <= 0 {
				// limit is of type time.Duration and the value is stored in nanoseconds. We convert it to minutes by dividing it by time.Minute 
				limit, _ := repo.GetResetTime(ip)
				remainingTime := time.Until(limit)
				return ShortURLData{}, fmt.Errorf("Rate limit is exceeded. Try again in %v minutes", remainingTime/time.Minute)
			}

			return ShortURLData{}, fmt.Errorf("Oops! Some error occurred in short url service")
		}

		// 2. check if the input is an actual URL
		if !govalidator.IsURL(url) {
			return ShortURLData{}, fmt.Errorf("Invalid URL")
		}

		// 3. Check for domain error
		if !utils.IsValidDomain(url) {
			return ShortURLData{}, fmt.Errorf("Invalid Domain")
		}

		// 4. enforce ssl for https

		url = utils.EnforceHTTP(url)



		/*
			After all checks have been passed, Create (or input from user) unique ID for url = domain + customShortCode.
			Check unique url is not already exists in DB as well.

		*/


		if CustomShortCode == "" {
			CustomShortCode = uuid.New().String()[:6]
		}

		if _,err = repo.GetURL(CustomShortCode); err != nil {
			return ShortURLData{}, fmt.Errorf("Custom short url is already in use")
		}


		id := uuid.New().String()[:6]
		
		// Save data in table ShortURL
		_ = repo.SaveURL(id, url, CustomShortCode, expiry, ip)

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
