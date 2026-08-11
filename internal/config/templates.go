package config

type TemplateInfo struct {
	Name          string           `yaml:"name"`
	Description   string           `yaml:"description"`
	Configuration map[string]Field `yaml:"configuration"`
}

type Field struct {
	Required    bool             `yaml:"required"`
	Type        string           `yaml:"type"`
	Default     any              `yaml:"default"`
	Description string           `yaml:"description"`
	Secret      bool             `yaml:"secret"`
	Values      []string         `yaml:"values"`
	Properties  map[string]Field `yaml:"properties"`
}
