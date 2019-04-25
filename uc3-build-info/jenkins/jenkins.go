package jenkins

import (
	"encoding/json"
	"fmt"
	. "github.com/CDLUC3/uc3-tools/mrt-build-info/maven"
	. "github.com/CDLUC3/uc3-tools/mrt-build-info/shared"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

// ------------------------------------------------------------
// Exported symbols

func IsParameterized(str string) bool {
	return paramSubRe.MatchString(str)
}

func Parameters(str string) []string {
	var parameters []string
	matches := paramSubRe.FindAllStringSubmatch(str, -1)
	for _, match := range matches {
		if len(match) != 2 { // should never happen
			panic(fmt.Errorf("invalid submatch: %#v", match))
		}
		parameters = append(parameters, match[1])
	}
	return parameters
}

// ------------------------------------------------------------
// Unexported symbols

func mapPomsToJobs(jobs []Job) (map[Pom]Job, []error) {
	//noinspection GoUnhandledErrorResult
	if Flags.Verbose {
		fmt.Fprintf(os.Stderr, "Retrieving POMs for %d jobs...", len(jobs))
		defer func() { fmt.Fprintln(os.Stderr, "Done.") }()
	}

	var errors []error
	jobsByPom := map[Pom]Job{}
	for _, j := range jobs {
		//noinspection GoUnhandledErrorResult
		if Flags.Verbose {
			fmt.Fprint(os.Stderr, ".")
		}
		poms, errs := j.POMs()
		if len(errs) > 0 {
			errors = append(errors, errs...)
		}
		if len(poms) == 0 {
			errors = append(errors, fmt.Errorf("no POMs found for job %v", j.Name()))
		}
		for _, p := range poms {
			jobsByPom[p] = j
		}
	}
	return jobsByPom, errors
}

var inTest = false
var client *http.Client
var apiUrlRelative = UrlMustParse("api/json?depth=1&pretty=true")
var apiUrlRegexp = regexp.MustCompile("/api/json(\\?.+)?$")

var configUrlRelative = UrlMustParse("config.xml")
var configUrlRegexp = regexp.MustCompile("/config.xml")

var paramSubRe = regexp.MustCompile("\\${([^}]+)}")

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

func isApiUrl(u *url.URL) bool {
	return apiUrlRegexp.MatchString(u.Path)
}

func toApiUrl(u *url.URL) *url.URL {
	if isApiUrl(u) {
		panic(fmt.Errorf("url '%v' is already an API URL", u))
	}
	return u.ResolveReference(apiUrlRelative)
}

func isConfigUrl(u *url.URL) bool {
	return configUrlRegexp.MatchString(u.Path)
}

func toConfigUrl(u *url.URL) *url.URL {
	if isConfigUrl(u) {
		panic(fmt.Errorf("url '%v' is already an CONFIG URL", u))
	}
	return u.ResolveReference(configUrlRelative)
}
