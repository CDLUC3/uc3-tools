package jenkins

type Parameter interface {
	Name() string
	Choices() []string
	Default() string
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

type defaultParamVal struct {
	Value string
}

