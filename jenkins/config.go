package jenkins

import (
	"github.com/beevik/etree"
	"github.com/dmolesUC3/mrt-build-info/shared"
	"net/url"
	"path/filepath"
	"regexp"
)

// ------------------------------------------------------------
// Exported symbols

type Config interface {
	URL() *url.URL
	Goals() string
	RootPOM() string
	BuildRoot() string
	MavenParameters() map[string]string
}

func ConfigFromURL(u *url.URL) (Config, error) {
	body, err := getBody(u)
	if err != nil {
		return nil, err
	}
	return ConfigFromBytes(body, u)
}

func ConfigFromBytes(data []byte, u *url.URL) (Config, error) {
	data = shared.HackXMLVersion(data)
	doc := etree.NewDocument()
	err := doc.ReadFromBytes(data)
	if err != nil {
		return nil, err
	}
	return &config{doc: doc, url: u}, nil
}

// ------------------------------------------------------------
// Unexported symbols

var propRe = regexp.MustCompile("-D([^=]+)=([^ ]+)")

type config struct {
	doc *etree.Document
	url *url.URL
}

func (c *config) URL() *url.URL {
	return c.url
}

func (c *config) Goals() string {
	goals := c.doc.FindElement("//goals")
	if goals == nil {
		return ""
	}
	return goals.Text()
}

func (c *config) RootPOM() string {
	rootPom := c.doc.FindElement("//rootPOM")
	if rootPom == nil {
		return ""
	}
	return rootPom.Text()
}

// Assumes we only ever build below the root POM
func (c *config) BuildRoot() string {
	rootPom := c.RootPOM()
	if rootPom == "" {
		return ""
	}
	return filepath.Dir(rootPom)
}

func (c *config) MavenParameters() map[string]string {
	paramToValue := map[string]string{}
	matches := propRe.FindAllStringSubmatch(c.Goals(), -1)
	for _, match := range matches {
		paramToValue[match[1]] = match[2]
	}
	return paramToValue
}