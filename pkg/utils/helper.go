package utils

import(
	"os"
	"strings"
)

/*
Function to enforce https
It changes http -> https in URL

Concepts used - slices
*/

func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "https://" + url
	}

	return url
}



/*
Function used to check for:
- Avoiding self-calls (service calling itself)
- Filtering internal domain requests
- Security validation
- Preventing redirects to same domain

Concepts used - strings and string operations
*/

func IsValidDomain(url string) bool {

	// Firsty, convert URL to fdqn (Fully qualified domain name)
	// i.e https://example1.com -> example1.com

	var newURL string;

	newURL = strings.Replace(url, "http://","",1) //"http://example.com/page/login" -> "example.com/page/login"
	newURL = strings.Replace(url, "https://","",1) // "https://example.com/page/login" -> "example.com/page/login"
	newURL = strings.Replace(url, "www.","",1) //"www.example.com/page/login" -> "example.com/page/login"
	newURL = strings.Split(url, "/")[0] // "example.com/page/login" -> example.com


	// Secondly, check if the given URL is the same as the domain configured in the environment variable, stop processing and return false.
	
	if newURL == os.Getenv("DOMAIN") {
		return false
	}

	return true
}