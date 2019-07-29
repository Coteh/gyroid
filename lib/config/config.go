package config

// Config describes a set of properties configurable by user
type Config struct {
	Clipboard bool
}

// ReadConfig parses config file and generates Config object containing the configurations
func ReadConfig(configFile string) (*Config, error) {

}

// WriteConfig writes configurations to file
func WriteConfig(config *Config, configFile string) error {

}
