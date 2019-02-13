package outputfmt

import (
	"fmt"
	"strings"
)

type Format interface {
	Name() string
	FieldSeparator() string
	InnerSeparator() string
	Prefix() string
	Suffix() string
}

func ToFormat(name string) (Format, error) {
	if format, ok := standardFormatNames[name]; ok {
		return &format, nil
	}
	return nil, fmt.Errorf("format name %#v should be one of: %v",
		name, strings.Join(StandardFormats(), ", "))
}

func StandardFormats() []string {
	var names []string
	for name := range standardFormatNames {
		names = append(names, name)
	}
	return names
}

type standardFormat int

const (
	TSV standardFormat = iota
	CSV
	Markdown
	Default = TSV
)

var standardFormatNames = map[string]standardFormat{
	"tsv": TSV,
	"csv": CSV,
	"md":  Markdown,
}

func (f standardFormat) FieldSeparator() string {
	switch f {
	case TSV:
		return "\t"
	case CSV:
		return ","
	case Markdown:
		return " | "
	default:
		return ""
	}
}

func (f standardFormat) InnerSeparator() string {
	switch f {
	case TSV:
		return ","
	case CSV:
		return ";"
	case Markdown:
		return " "
	default:
		return ""
	}
}

func (f standardFormat) Prefix() string {
	if f == Markdown {
		return "| "
	}
	return ""
}

func (f standardFormat) Suffix() string {
	if f == Markdown {
		return " |"
	}
	return ""
}

func (f standardFormat) Name() string {
	for k, v := range standardFormatNames {
		if v == f {
			return k
		}
	}
	return ""
}
