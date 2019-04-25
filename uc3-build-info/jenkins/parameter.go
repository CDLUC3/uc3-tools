package jenkins

import (
	"sort"
	"strings"
)

type Parameter interface {
	Name() string
	Choices() []string
	Default() string
	Parameterize(str string) []string
}

type parameterDefinition struct {
	Class string `json:"_class"`
	Name_ string `json:"name"`
	Choices_ []string `json:"choices"`
	DefaultParameterValue *defaultParamVal
}

func (p *parameterDefinition) Name() string {
	return p.Name_
}

func (p *parameterDefinition) Choices() []string {
	return p.Choices_
}

func (p *parameterDefinition) Default() string {
	def := p.DefaultParameterValue
	if def != nil {
		return def.Value
	}
	return ""
}

func (p *parameterDefinition) Parameterize(str string) []string {
	paramSub := "${" + p.Name() + "}"
	if strings.Contains(str, paramSub) {
		choices := p.Choices()
		sort.Strings(choices)
		parameterized := make([]string, len(choices))
		for i, c := range choices {
			parameterized[i] = strings.ReplaceAll(str, paramSub, c)
		}
		return parameterized
	}
	return []string{str}
}

type defaultParamVal struct {
	Value string
}

