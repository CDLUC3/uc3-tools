package jenkins

import (
	"fmt"
	"net/url"
	"regexp"
)

// ------------------------------------------------------------
// Unexported symbols

func urlMustParse(urlStr string) *url.URL {
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	return u
}

var apiUrlRelative = urlMustParse("api/json?depth=1&pretty=true")
var apiUrlRegexp = regexp.MustCompile("/api/json(\\?.+)?$")

func toApiUrl(u *url.URL) *url.URL {
	if apiUrlRegexp.MatchString(u.Path) {
		panic(fmt.Errorf("url '%v' is already an API URL", u))
	}
	return u.ResolveReference(apiUrlRelative)
}