package appconfig

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type AiChatConfig struct {
	Model              string   `yaml:"model"`
	Temperature        float32  `yaml:"temperature"`
	SystemMessage      string   `yaml:"systemMessage"`
	AvailableFunctions []string `yaml:"availableFunctions"`
	TmpHttpPort        int      `yaml:"tmpHttpPort"`
}

type FunctionConfig struct {
	Url             string `yaml:"url"`
	Description     string `yaml:"description"`
	RequestTemplate string `yaml:"requestTemplate"`
}

type AppConfig struct {
	AiChatCfg   map[string]*AiChatConfig   `yaml:"chats"`
	FunctionCfg map[string]*FunctionConfig `yaml:"functions"`
}

var AppCfg *AppConfig

func LoadConfig(path string) error {
	log.Println("Load config file:", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read file %s: %w", path, err)
	}
	if err := yaml.Unmarshal(data, &AppCfg); err != nil {
		return fmt.Errorf("YAML parsing error: %w", err)
	}
	log.Println("Config file loaded")
	return nil
}
