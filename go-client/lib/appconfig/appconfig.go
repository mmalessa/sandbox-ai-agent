package appconfig

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	AiChatCfg map[string]*AiChatConfig `yaml:"chats"`
}

type AiChatConfig struct {
	Model              string   `yaml:"model"`
	Temperature        float32  `yaml:"temperature"`
	SystemMessage      string   `yaml:"systemMessage"`
	AvailableFunctions []string `yaml:"availableFunctions"`
	TmpHttpPort        int      `yaml:"tmpHttpPort"`
}

var AiChatCfg *AiChatConfig
var AiChatCfgs []*AiChatConfig

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
	if err := validateConfig(); err != nil {
		return err
	}
	log.Println("Config file loaded")
	return nil
}

func validateConfig() error {
	// defaultChat := AppCfg.DefaultChat
	// if _, ok := AppCfg.AiChatCfg[defaultChat]; !ok {
	// 	return fmt.Errorf("configuration for chat \"%s\" not found", defaultChat)
	// }
	return nil
}
