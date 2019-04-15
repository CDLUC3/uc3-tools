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
	Repository() Repository
	Path() string
	GetContent() ([]byte, error)
}

type entry struct {
	path  string
	sha1  string
	eType EntryType
	size  int
	url   *url.URL
	repository *repository
}

func (e *entry) Repository() Repository {
	return e.repository
}

func (e *entry) Path() string {
	return e.path
}

func (e *entry) GetContent() ([]byte, error) {
	if e.eType != Blob {
		return nil, fmt.Errorf("can't get content of %v entry", e.eType)
	}
	u := e.url
	//noinspection GoBoolExpressions
	if inTest && !strings.HasPrefix(u.Host, "127.0.0.1") {
		return nil, fmt.Errorf("no real URLs in test!: %v", u)
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", contentTypeRaw)

	httpClient := e.repository.HttpClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(bytes) != e.size {
		return nil, fmt.Errorf("expected %d bytes, got %d", e.size, len(bytes))
	}
	return bytes, nil
}