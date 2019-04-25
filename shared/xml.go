package shared

import (
	"regexp"
)

var xmlVersionRe = regexp.MustCompile("<\\?xml version=['\"]1\\.[10]['\"]([^>]+)>")
var xmlVersion1_0 = []byte("<?xml version='1.0'$1>")

func HackXMLVersion(data []byte) []byte {
	// TODO: other XML 1.1 differences?
	//   - see https://www.w3.org/TR/xml11/#sec-xml11
	return xmlVersionRe.ReplaceAll(data, xmlVersion1_0)
}
