package jenkins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// ------------------------------------------------------------
// Unexported symbols

var inTest = false

var client *http.Client

var apiUrlRelative = urlMustParse("api/json?depth=1&pretty=true")
var apiUrlRegexp = regexp.MustCompile("/api/json(\\?.+)?$")

func getBody(u *url.URL) ([]byte, error) {
	//noinspection GoBoolExpressions
	if inTest && !strings.HasPrefix(u.Host, "127.0.0.1") {
		return nil, fmt.Errorf("no real URLs in test!: %v", u)
	}

	if client == nil {
		client = &http.Client{Timeout: time.Second * 30}
	}
	resp, err := client.Get(u.String())
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func unmarshal(u *url.URL, target interface{}) (err error) {
	body, err := getBody(u)
	if err == nil {
		err = json.Unmarshal(body, target)
	}
	return
}

func urlMustParse(urlStr string) *url.URL {
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	return u
}

func toApiUrl(u *url.URL) *url.URL {
	if apiUrlRegexp.MatchString(u.Path) {
		panic(fmt.Errorf("url '%v' is already an API URL", u))
	}
	return u.ResolveReference(apiUrlRelative)
}
