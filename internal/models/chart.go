package models

type MiniHelmChart struct {
	Name         string       `yaml:"name,omitempty"`
	Version      string       `yaml:"version,omitempty"`
	Path         string       `yaml:"path,omitempty"`
	Dependencies []Dependency `yaml:"dependencies,omitempty"`
	Repository   string       `yaml:"repository,omitempty"`
}

type Dependency struct {
	Name    string `yaml:"name,omitempty"`
	Version string `yaml:"version,omitempty"`

	Repository   string   `yaml:"repository,omitempty"`
	Condition    string   `yaml:"condition,omitempty"`
	Tags         []string `yaml:"tags,omitempty"`
	Enabled      bool     `yaml:"enabled,omitempty"`
	ImportValues []any    `yaml:"import-values,omitempty"`
	Alias        string   `yaml:"alias,omitempty"`
}
