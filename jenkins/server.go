package jenkins

import (
	"fmt"
	"net/url"
)

// Server represents a Jenkins server
type Server interface {
}

func DefaultServer() Server {
	server, err := ServerFromUrl("http://builds.cdlib.org/")
	if err == nil {
		return server
	}
	panic(err)
}

func ServerFromUrl(urlStr string) (Server, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if !u.IsAbs() {
		return nil, fmt.Errorf("server URL '%v' is not absolute", urlStr)
	}
	return &server{serverUrl: u}, nil
}

// ------------------------------------------------------------
// Unexported symbols

type server struct {
	serverUrl *url.URL
}

func (s *server) apiRoot() *url.URL {
	return toApiUrl(s.serverUrl)
}