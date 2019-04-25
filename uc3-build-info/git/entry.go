package git

import (
	"fmt"
	"github.com/CDLUC3/uc3-tools/mrt-build-info/shared"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const contentTypeRaw = "application/vnd.github.v3.raw"

type Entry interface {
	Repository() Repository
	SHA1() SHA1
	Path() string
	GetContent() ([]byte, error)
	URL() *url.URL
}

func WebUrlForEntry(e Entry) *url.URL {
	repo := e.Repository()
	sha1 := repo.SHA1()
	if !FullSHA {
		sha1 = sha1[0:12]
	}
	u := fmt.Sprintf("http://github.com/%v/%v/blob/%v/%v", repo.Owner(), repo.Name(), sha1, e.Path())
	return shared.UrlMustParse(u)
}

type entry struct {
	path       string
	sha1       SHA1
	eType      EntryType
	size       int
	url        *url.URL
	repository *repository

	content []byte
}

func (e *entry) Repository() Repository {
	return e.repository
}

func (e *entry) SHA1() SHA1 {
	return e.sha1
}

func (e *entry) Path() string {
	return e.path
}

func (e *entry) URL() *url.URL {
	return e.url
}

func (e *entry) GetContent() ([]byte, error) {
	if e.eType != Blob {
		return nil, fmt.Errorf("can't get content of %v entry", e.eType)
	}

	if e.content == nil {
		u := e.URL()
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
		e.content = bytes
	}
	return e.content, nil
}
