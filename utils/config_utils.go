package utils

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `yaml:"server"`
	Static HtmlConfig   `yaml:"html"`
}
type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
type HtmlConfig struct {
	Dirpath string `yaml:"path"`
	Index   string `yaml:"index"`
}

func LoadConfig() *Config {
	var config Config
	configPath := filepath.Join("conf", "gohttp.yaml")
	confData, readErr := os.ReadFile(configPath)
	if readErr != nil {
		log.Fatalln(readErr)
	}
	if unmarshalErr := yaml.Unmarshal(confData, &config); unmarshalErr != nil {
		log.Fatalln(readErr)
	}
	return &config
}
