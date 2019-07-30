package config

import (
	"errors"
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

const NO_CONFIG_WRITE_ERROR = "No config to write"
const NO_BYTES_WRITTEN_ERROR = "No bytes written to file"

// Config describes a set of properties configurable by user
type Config struct {
	Clipboard bool
}

// ReadConfig parses config file and generates Config object containing the configurations
func ReadConfig(reader io.Reader) (*Config, error) {
	configBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	config := new(Config)

	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// WriteConfig writes configurations to file
func WriteConfig(config *Config, writer io.Writer) error {
	if config == nil {
		return errors.New(NO_CONFIG_WRITE_ERROR)
	}

	configBytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	bytesWritten, err := writer.Write(configBytes)
	if err != nil {
		return err
	}
	if bytesWritten == 0 {
		return errors.New(NO_BYTES_WRITTEN_ERROR)
	}

	return nil
}
