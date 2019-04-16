package jenkins

import (
	"fmt"
	"github.com/beevik/etree"
	"net/url"
)

type Config interface {
	Goals() (string, error)
}

type config struct {
	url *url.URL
	doc *etree.Document
}

func (c *config) Goals() (string, error) {
	doc, err := c.document()
	if err != nil {
		return "", err
	}
	goals := doc.FindElement("//goals")
	if goals != nil {
		return "", fmt.Errorf("<goals> not found in %v", c.url)
	}
	return goals.Text(), nil
}

func (c *config) document() (*etree.Document, error) {
	if c.doc == nil {
		body, err := getBody(c.url)
		if err != nil {
			return nil, err
		}

		doc := etree.NewDocument()
		err = doc.ReadFromBytes(body)
		if err != nil {
			return nil, err
		}
		c.doc = doc
	}
	return c.doc, nil
}
