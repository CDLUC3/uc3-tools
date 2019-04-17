package jenkins

import (
	"fmt"
	"github.com/beevik/etree"
	"github.com/dmolesUC3/mrt-build-info/misc"
	"net/url"
)

// ------------------------------------------------------------
// Exported symbols

type Config interface {
	URL() *url.URL
	Goals() (string, error)
}

func ConfigFromURL(u *url.URL) (Config, error) {
	body, err := getBody(u)
	if err != nil {
		return nil, err
	}
	return ConfigFromBytes(body, u)
}

func ConfigFromBytes(data []byte, u *url.URL) (Config, error) {
	data = misc.HackXMLVersion(data)
	doc := etree.NewDocument()
	err := doc.ReadFromBytes(data)
	if err != nil {
		return nil, err
	}
	return &config{doc: doc, url: u}, nil
}

// ------------------------------------------------------------
// Unexported symbols

type config struct {
	doc *etree.Document
	url *url.URL
}

func (c *config) URL() *url.URL {
	return c.url
}

func (c *config) Goals() (string, error) {
	goals := c.doc.FindElement("//goals")
	if goals == nil {
		return "", fmt.Errorf("<goals> not found in %v", c.url)
	}
	return goals.Text(), nil
}
