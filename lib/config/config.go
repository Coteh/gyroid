package config

import (
	"errors"
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// NoConfigWriteError is an error message for no config file to write to
const NoConfigWriteError = "No config to write"

// NoBytesWrittenError is an error message for no bytes written
const NoBytesWrittenError = "No bytes written to file"

// Config describes a set of properties configurable by user
type Config struct {
	Clipboard bool
	RefineNew string `yaml:"refine-new"`
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
		return errors.New(NoConfigWriteError)
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
		return errors.New(NoBytesWrittenError)
	}

	return nil
}
