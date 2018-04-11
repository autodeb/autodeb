package app

// Config configures the application
type Config struct {
	TemplatesDirectory string
	DataDirectory      string
}

// NewConfig creates an application config
func NewConfig(templatesDirectory, dataDirectory string) *Config {
	cfg := Config{
		TemplatesDirectory: templatesDirectory,
		DataDirectory:      dataDirectory,
	}
	return &cfg
}
