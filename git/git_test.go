package git

import (
	"fmt"
	"github.com/dmolesUC3/mrt-build-info/misc"
	. "gopkg.in/check.v1"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	inTest = true
	TestingT(t)
}

var _ = Suite(&GitSuite{})

type GitSuite struct {
	server    *httptest.Server
	serverUrl *url.URL
}

var data = map[string]string{
	"/repos/CDLUC3/mrt-store/git/trees/af174ac555758a1c639a7a3da39e022d9fdbf3a6": "testdata/tree.json",
	"/repos/CDLUC3/mrt-store/git/blobs/d3ae87b904091324ea42af37581d92ee67415143": "testdata/pom.xml",
}

// HandleRequest is a hack to rewrite the test data to replace real URLs with
// URLs from the local httptest.Server
func (s *GitSuite) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if file, ok := data[r.URL.Path]; ok {
		if strings.HasSuffix(file, ".json") || r.Header.Get("Accept") == contentTypeRaw {
			body, err := ioutil.ReadFile(file)
			if err == nil {
				// note that URL.Host includes port
				body = []byte(strings.ReplaceAll(string(body), "api.github.com", s.serverUrl.Host))
				n, err := w.Write(body)
				if err != nil {
					panic(err)
				}
				if n != len(body) {
					panic(fmt.Errorf("wrong number of bytes: wrote %d, expected %d", n, len(body)))
				}
				return
			}
		}
	}
	w.WriteHeader(500)
}

// SetUpTest sets up the local httptest.Server -- note that it can't be for the whole
// suite since httptest.Server is designed to only live for one testing.T test
func (s *GitSuite) SetUpTest(c *C) {
	s.server = httptest.NewServer(http.HandlerFunc(s.HandleRequest))
	s.serverUrl = misc.UrlMustParse(s.server.URL)
}

func (s *GitSuite) MockClient() *http.Client {
	return &http.Client{Transport: s}
}

// RoundTrip Implements http.RoundTripper
func (s *GitSuite) RoundTrip(req *http.Request) (*http.Response, error) {
	// TODO: just use this instead of httptest.Server?
	req.URL.Host = s.serverUrl.Host
	if req.URL.Scheme == "https" {
		req.URL.Scheme = "http"
	}
	return http.DefaultTransport.RoundTrip(req)
}

func (s *GitSuite) TestEntries(c *C) {
	repo := NewRepository("CDLUC3", "mrt-store", "af174ac555758a1c639a7a3da39e022d9fdbf3a6", "").(*repository)
	repo.httpClient = s.MockClient()

	entries, err := repo.Find("pom.xml", Blob)
	c.Assert(err, IsNil)
	c.Assert(len(entries), Equals, 3)

	var entry Entry
	for _, e := range entries {
		if e.Path() == "store-src/pom.xml" {
			entry = e
			break
		}
	}
	c.Assert(entry, NotNil)

	content, err := entry.GetContent()
	c.Assert(err, IsNil)

	expectedBytes, _ := ioutil.ReadFile("testdata/pom.xml")
	expected := string(expectedBytes)

	c.Assert(content, Equals, expected)
}
