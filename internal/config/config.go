package config

type Config struct {
	Project  string         `yaml:"project"`
	Compose  ComposeConfig  `yaml:"compose"`
	Docker   DockerConfig   `yaml:"docker"`
	Database DatabaseConfig `yaml:"database"`
}

type ComposeConfig struct {
	File string `yaml:"file"`
}

type DockerConfig struct {
	Project string `yaml:"project"`
}

type DatabaseConfig struct {
	User string `yaml:"user"`
	Host string `yaml:"host"`
	Pass string `yaml:"pass"`
	Port int    `yaml:"port"`
	Name string `yaml:"name"`
}
