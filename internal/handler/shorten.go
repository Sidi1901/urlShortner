package handler

import(
   "time"
   "strconv"
   "os"
   "fmt"
   "github.com/Sidi1901/urlShortner/internal/repository"
   repo "github.com/Sidi1901/urlShortner/pkg/utils/utilities"
   "github.com/asaskevich/govalidator"
   "github.com/google/uuid"
)

type request struct {
	URL	          string             `json:"url"`
	CustomShort   string             `json:"short"`
	Expiry        time.Duration      `json:"expiry"`
}

type response struct {
	URL              string          `json:"url"`
	CustomShort      string          `json:"shorturl"`
	Expiry           time.Duration   `json:"expiry"`
	XRateRemaining   int             `json:"rate_limit"`
	XRateLimitReset  time.Duration   `json:"rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {
	
	body := new(request)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"cannot parse json"})
	}



	// 1. Implement rate limiting

	// Try to get the value of the IP address from the database. If the value does not exist, it means that it's the first time the IP address is making a request, so we set the rate limit for that IP address. If the value exists, we check if it's less than or equal to 0, which means that the rate limit has been exceeded. If the rate limit is exceeded, we return a 503 Service Unavailable status with a message indicating that the rate limit has been exceeded and when it will reset.

	_, err := repo.GetRemainingQuota(c.IP())

	if err == nil{
		// Rate limit exists for the IP address. Check if it's exceeding the quota. 
	
		// Check if val <=0 indicating expired key in db which was initially set to 30 minues.
		val, _ := repo.GetRemainingQuota(c.IP())
		
		if val <= 0 {
			
		}





		id := uuid.New().String()[:6]

		var customShort string

		if body.CustomShort == "" {
			customShort = uuid.New().String()[:6]
		} else {
			customShort = body.CustomShort
		}
		_ = repo.SaveURL(id, customShort, body.URL c.IP(), os.Getenv("API_QUOTA"),30*60*time.Second).Err()

	} else {

		// Check if val <=0 indicating expired key in redis which was initially set to 30 minues.
		val, _ := r2.Get(database.Ctx, c.IP()).Result()
		valInt, _ := strconv.Atoi(val)

		if valInt <= 0 {
			// limit is of type time.Duration and the value is stored in nanoseconds. We convert it to minutes by dividing it by time.Minute 
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()

			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": "Rate limit is exceeded",
				"rate_limit_reset": limit / time.Minute,
			})
		}
	}


	// 2. check if the input is an actual URL
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}

	// 3. Check for domain error

	if !helpers.IsValidDomain(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Invalid domain"})
	}

	// 4. enforce ssl for https

	body.URL = helpers.EnforceHTTP(body.URL)


	/*
		After all checks have been passed, Create (or input from user) unique ID for url = domain + id.
		Check unique url is not already exists in DB as well.

	*/

	var id string

	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.CreateClient(0)
	defer r.Close()

	val, _ := r.Get(database.Ctx, id).Result()

	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Custom short url is already in use",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	fmt.Println(body.URL)
	fmt.Println(id)

	if err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to connect to the database",
		})
	}

	resp := response{
		URL:              body.URL,
		CustomShort:      "",
		Expiry:           body.Expiry,
		XRateRemaining:   10,
		XRateLimitReset:  30,
	}

	r2.Decr(database.Ctx, c.IP())

	val, _ = r2.Get(database.Ctx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()

	resp.XRateLimitReset = ttl / time.Nanosecond

	resp.CustomShort = os.Getenv("Domain") + "/" + id

	return c.Status(fiber.StatusOK).JSON(resp)


}

	// The rate limit is set to the value of the environment variable "API_QUOTA" and it expires after 30 minutes (30*60 seconds).