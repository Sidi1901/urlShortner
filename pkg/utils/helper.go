package utils

import (
	"os"
	"regexp"
	"strings"
)

var (
	hasLetter  = regexp.MustCompile(`[A-Za-z]`)
	hasNumber  = regexp.MustCompile(`\d`)
	hasSpecial = regexp.MustCompile(`[@$!%*#?&]`)
)

func IsValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	if !hasLetter.MatchString(password) {
		return false
	}
	if !hasNumber.MatchString(password) {
		return false
	}
	if !hasSpecial.MatchString(password) {
		return false
	}
	return true
}

func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "https://" + url
	}

	return url
}

func IsValidDomain(url string) bool {

	// Firsty, convert URL to fdqn (Fully qualified domain name)
	// i.e https://example1.com -> example1.com

	var newURL string

	newURL = strings.Replace(url, "http://", "", 1)     //"http://example.com/page/login" -> "example.com/page/login"
	newURL = strings.Replace(newURL, "https://", "", 1) // "https://example.com/page/login" -> "example.com/page/login"
	newURL = strings.Replace(newURL, "www.", "", 1)     //"www.example.com/page/login" -> "example.com/page/login"
	newURL = strings.Split(newURL, "/")[0]              // "example.com/page/login" -> example.com

	// Secondly, check if the given URL is the same as the domain configured in the environment variable, stop processing and return false.

	if newURL == os.Getenv("DOMAIN") {
		return false
	}

	return true
}
