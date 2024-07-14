package utils

import (
	"log"
	"os"
	"path/filepath"
	"reflect"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig    `yaml:"server,omitempty"`
	Static  HtmlConfig      `yaml:"html,omitempty"`
	Custom  []CustomConfig  `yaml:"custom,omitempty"`
	Proxy   []ProxyConfig   `yaml:"proxy,omitempty"`
	Backend []BackendConfig `yaml:"backend,omitempty"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type HtmlConfig struct {
	Dirpath string `yaml:"path"`
	Index   string `yaml:"index"`
}

type LoggerConfig struct {
}

type CustomConfig struct {
	Urlpath  string `yaml:"url"`
	Filepath string `yaml:"file"`
}

type ProxyConfig struct {
	PathPrefix string `yaml:"prefix"`
	TargetURL  string `yaml:"target"`
}

type BackendConfig struct {
	BackendURL string `yaml:"url"`
}

func DefaultServer() *ServerConfig {
	return &ServerConfig{
		Host: "0.0.0.0",
		Port: "80",
	}
}

func DefaultHtml() *HtmlConfig {
	return &HtmlConfig{
		Dirpath: "html",
		Index:   "index.html",
	}
}

func CoverConfig(c *Config) {
	if reflect.DeepEqual(c.Server, ServerConfig{}) {
		c.Server = *DefaultServer()
	}

	if reflect.DeepEqual(c.Static, HtmlConfig{}) {
		c.Static = *DefaultHtml()
	}
}

func LoadConfig() *Config {
	var config Config
	configPath := filepath.Join("conf", "gohttp.yaml")
	confData, readErr := os.ReadFile(configPath)
	if readErr != nil {
		log.Fatalln(readErr)
	}
	if unmarshalErr := yaml.Unmarshal(confData, &config); unmarshalErr != nil {
		log.Fatalln(unmarshalErr)
	}
	CoverConfig(&config)
	return &config
}
