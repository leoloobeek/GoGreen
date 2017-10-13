package lib

import (
	"io/ioutil"
	"net/http"
)

// GenerateHttpKey reaches out to url, hashes the body and returns
// the sha512 hash (first 32 chars)
func GenerateHttpKey(url, userAgent string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	httpKey := GenerateSHA512(string(result))
	return httpKey[:32], nil

}
