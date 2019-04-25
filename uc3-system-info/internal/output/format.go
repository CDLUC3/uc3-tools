package output

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Format interface {
	Name() string
	FieldSeparator() string
	InnerSeparator() string
	Prefix() string
	Suffix() string
	SprintTitle(title string) string
	SprintExample(example string) string
	SprintHeader(headerFields ...string) string
	Sprint(fields ...interface{}) (string, error)
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

func (f standardFormat) SprintTitle(title string) string {
	if f == Markdown {
		return fmt.Sprintf("### %v\n", title)
	}
	return title + "\n"
}

func (f standardFormat) SprintExample(example string) string {
	if f == Markdown {
		return fmt.Sprintf("`%v`", example)
	}
	return example
}

func (f standardFormat) SprintHeader(headerFields ...string) string {
	var sb strings.Builder
	sb.WriteString(f.Prefix())
	for i, hf := range headerFields {
		sb.WriteString(hf)
		if i+1 < len(headerFields) {
			sb.WriteString(f.FieldSeparator())
		}
	}
	sb.WriteString(f.Suffix())
	sb.WriteString("\n")

	if f == Markdown {
		sb.WriteString(f.Prefix())
		for i := 0; i < len(headerFields); i++ {
			sb.WriteString(":---")
			if i+1 < len(headerFields) {
				sb.WriteString(f.FieldSeparator())
			}
		}
		sb.WriteString(f.Suffix())
		sb.WriteString("\n")
	}
	return sb.String()
}

func (f standardFormat) Sprint(fields ...interface{}) (string, error) {
	var sb strings.Builder
	sb.WriteString(f.Prefix())
	for i, field := range fields {
		switch v := field.(type) {
		case string:
			sb.WriteString(v)
		case int64:
			sb.WriteString(strconv.FormatInt(v, 10))
		case fmt.Stringer:
			sb.WriteString(v.String())
		case []string:
			for j, sf := range v {
				sb.WriteString(sf)
				if j+1 < len(v) {
					sb.WriteString(f.InnerSeparator())
				}
			}
		default:
			return "", fmt.Errorf("don't know how to format %v %#v", reflect.TypeOf(field), field)
		}
		if i+1 < len(fields) {
			sb.WriteString(f.FieldSeparator())
		}
	}
	sb.WriteString(f.Suffix())
	return sb.String(), nil
}
