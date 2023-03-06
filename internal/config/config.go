package config

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

var k = koanf.New(".")

// Config contains the main snaprd config
type Config struct {
	Schedule string   `koanf:"schedule"`
	Snapraid Snapraid `koanf:"snapraid"`
	Scrub    Scrub    `koanf:"scrub"`
}

// Snapraid contains the snapraid specific config
type Snapraid struct {
	Executable      string `koanf:"executable"`
	ConfigPath      string `koanf:"configPath"`
	DeleteThreshold int    `koanf:"deleteThreshold"`
	Touch           bool   `koanf:"touch"`
}

// Scrub contains the snapraid scrub specific config
type Scrub struct {
	Enabled    bool `koanf:"enabled"`
	Percentage int  `koanf:"percentage"`
	OlderThan  int  `koanf:"olderThan"`
}

// Parse loads configuration file and unmarshal's to struct
func Parse(configFile string) (*Config, error) {
	err := k.Load(file.Provider(configFile), yaml.Parser())
	if err != nil {
		return nil, err
	}
	var conf Config
	if err := k.Unmarshal("", &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
