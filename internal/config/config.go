package config

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

var k = koanf.New(".")

type Config struct {
	Schedule string   `koanf:"schedule"`
	Snapraid Snapraid `koanf:"snapraid"`
	Scrub    Scrub    `koanf:"scrub"`
}

type Snapraid struct {
	Executable      string `koanf:"executable"`
	Config          string `koanf:"config"`
	DeleteThreshold int    `koanf:"deleteThreshold"`
	Touch           bool   `koanf:"touch"`
}

type Scrub struct {
	Enabled    bool `koanf:"enabled"`
	Percentage int  `koanf:"percentage"`
	OlderThan  int  `koanf:"olderThan"`
}

// Load configuration file and unmarshal to struct
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
