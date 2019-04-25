package shared

import (
	"net/url"
)

func UrlMustParse(urlStr string) *url.URL {
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	return u
}

