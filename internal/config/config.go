package config

type Config struct {
	Project string        `yaml: "project"`
	Compose ComposeConfig `yaml: "compose"`
	Docker  DockerConfig  `yaml: "docker"`
}

type ComposeConfig struct {
	File string `yaml:"file"`
}

type DockerConfig struct {
	Project string `yaml:"project"`
}
