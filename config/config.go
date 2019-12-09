package config

// Config ...
type Config struct {
	LogLevel   string `toml:"log_level"`
	BindAddr   string `toml:"bind_address"`
	ConfigPath string `toml:"path_config"`
}

//NewConfig ...
func NewConfig() *Config {
	return &Config{
		LogLevel:   "",
		BindAddr:   "",
		ConfigPath: "",
	}
}
