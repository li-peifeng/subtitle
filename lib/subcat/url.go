package subcat

import (
	"net/url"
)

const (
	baseURL = "https://subtitlecat.com/"
)

func buildURL(baseURL string, path string, params map[string]string) string {
	u, _ := url.Parse(baseURL)

	u.Path = path

	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	return u.String()
}
