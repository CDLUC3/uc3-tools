package output

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
	FormatHeader(header []string) string
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
		return "<br/>"
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

func (f standardFormat) FormatHeader(header []string) string {
	headerLine := strings.Join(header, f.FieldSeparator())
	headerLine = fmt.Sprintf("%v%v%v\n", f.Prefix(), headerLine, f.Suffix())
	if f == Markdown {
		var sb strings.Builder
		sb.WriteString(headerLine)
		sb.WriteString(f.Prefix())
		for i := 0; i < len(header); i++ {
			sb.WriteString(":---")
			if i + 1 < len(header) {
				sb.WriteString(f.FieldSeparator())
			}
		}
		sb.WriteString(f.Suffix())
		sb.WriteString("\n")
		return sb.String()
	}
	return headerLine
}