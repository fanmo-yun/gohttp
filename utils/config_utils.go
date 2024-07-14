package utils

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig    `yaml:"server,omitempty"`
	Static  HtmlConfig      `yaml:"html,omitempty"`
	Logger  LoggerConfig    `yaml:"logger,omitempty"`
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
	Try     bool
}

type LoggerConfig struct {
	Out   string `yaml:"out"`
	Level string `yaml:"level"`
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
		Try:     false,
	}
}

func DefaultLogger() *LoggerConfig {
	return &LoggerConfig{
		Out:   "stdout",
		Level: "info",
	}
}

func CoverConfig(c *Config) {
	if reflect.DeepEqual(c.Server, ServerConfig{}) {
		c.Server = *DefaultServer()
	} else {
		if c.Server.Host == "" {
			c.Server.Host = "0.0.0.0"
		} else if c.Server.Port == "" {
			c.Server.Port = "80"
		}
	}

	if reflect.DeepEqual(c.Static, HtmlConfig{}) {
		c.Static = *DefaultHtml()
	} else {
		if c.Static.Dirpath == "" {
			c.Static.Dirpath = "html"
		} else if c.Static.Index == "" {
			c.Static.Index = "index.html"
		} else {
			s := strings.Split(c.Static.Index, " ")
			if strings.ToLower(s[0]) == "try" {
				c.Static.Try = true
				if len(s) > 1 {
					c.Static.Index = s[1]
				} else {
					panic("配置错误: 'Index' 字段格式错误，缺少文件名")
				}
			}
		}
	}

	if reflect.DeepEqual(c.Logger, LoggerConfig{}) {
		c.Logger = *DefaultLogger()
	} else {
		if c.Logger.Out == "" {
			c.Logger.Out = "stdout"
		} else if c.Logger.Level == "" {
			c.Logger.Level = "info"
		}
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
