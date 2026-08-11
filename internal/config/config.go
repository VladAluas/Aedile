package config

type Config struct {
	Project   ProjectConfig   `yaml:"project"`
	Compose   ComposeConfig   `yaml:"compose"`
	Docker    DockerConfig    `yaml:"docker"`
	Metadata  MetadataConfig  `yaml:"metadata"`
	Database  DatabaseConfig  `yaml:"database"`
	Storage   StorageConfig   `yaml:"storage"`
	Liquibase LiquibaseConfig `yaml:"liquibase"`
}

type ProjectConfig struct {
	Name        string `yaml:"name"`
	Environment string `yaml:"environment"`
}

type ComposeConfig struct {
	File string `yaml:"file"`
}

type DockerConfig struct {
	Project string `yaml:"project"`
}

type MetadataConfig struct {
	Engine string `yaml:"engine"`
	Image  string `yaml:"image"`
	User   string `yaml:"user"`
	Host   string `yaml:"host"`
	Pass   string `yaml:"pass"`
	Port   int    `yaml:"port"`
	Name   string `yaml:"name"`
}

type DatabaseConfig struct {
	Engine string `yaml:"engine"`
	Image  string `yaml:"image"`
	User   string `yaml:"user"`
	Host   string `yaml:"host"`
	Pass   string `yaml:"pass"`
	Port   int    `yaml:"port"`
	Name   string `yaml:"name"`
}

type StorageConfig struct {
	Image       string `yaml:"image"`
	User        string `yaml:"user"`
	Pass        string `yaml:"pass"`
	APIPort     int    `yaml:"api_port"`
	WebPort     int    `yaml:"web_port"`
	Container   string `yaml:"container"`
	Command     string `yaml:"command"`
	StorageData string `yaml:"storage"`
}

type LiquibaseConfig struct {
	Image      string `yaml:"image"`
	Name       string `yaml:"name"`
	WorkingDir string `yaml:"working_dir"`
	Changelog  string `yaml:"changelog"`
	Properties string `yaml:"properties"`
	Dependency string `yaml:"dependency"`
	Lib        string `yaml:"lib"`
}
