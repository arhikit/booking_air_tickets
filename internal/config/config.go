package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

var cfg Config

type Config struct {
	BasePath string `yaml:"base_path"`
	Addr     string `yaml:"addr"`
	DB       struct {
		Postgresql string `yaml:"postgresql"`
	} `yaml:"db"`
}

func InitConfig(args []string) (*Config, error) {
	var configPath string

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.StringVar(&configPath, "c", "./cmd/app/config.yaml", "set path to config")
	err := flags.Parse(args[1:])
	if err != nil {
		return nil, err
	}

	clean := filepath.Clean(configPath)

	file, err := os.Open(clean)
	if err != nil {
		return nil, fmt.Errorf("fail to open config file in path \"%s\" with error %w", configPath, err)
	}

	err = yaml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("fail to parse config %w", err)
	}

	return &cfg, nil
}
