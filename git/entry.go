package git

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const contentTypeRaw = "application/vnd.github.VERSION.raw"

type Entry interface {
	Path() string
	GetContent() (string, error)
}

type entry struct {
	path  string
	sha1  string
	eType EntryType
	size  int
	url   *url.URL
	repository *repository
}

func (e *entry) Path() string {
	return e.path
}

func (e *entry) GetContent() (string, error) {
	if e.eType != Blob {
		return "", fmt.Errorf("can't get content of %v entry", e.eType)
	}
	u := e.url
	//noinspection GoBoolExpressions
	if inTest && !strings.HasPrefix(u.Host, "127.0.0.1") {
		return "", fmt.Errorf("no real URLs in test!: %v", u)
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", contentTypeRaw)

	httpClient := e.repository.HttpClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if len(bytes) != e.size {
		return "", fmt.Errorf("expected %d bytes, got %d", e.size, len(bytes))
	}
	return string(bytes), nil
}