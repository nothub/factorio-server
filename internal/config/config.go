package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
)

var Loaded Configuration

type Configuration struct {
	FactorioUser  string   `yaml:"factorio-user"`
	FactorioToken string   `yaml:"factorio-token"`
	WebhookChat   string   `yaml:"webhook-chat"`
	WebhookStatus string   `yaml:"webhook-status"`
	Mods          []string `yaml:"mods"`
}

func Load() {

	ReadFile("/etc/factorio-server/config.yaml")
	ReadFile(filepath.Join(os.Getenv("HOME"), ".config/factorio-server/config.yaml"))
	ReadFile(filepath.Join(os.Getenv("PWD"), "server/config.yaml"))

	Clean()
	Validate()

}

func ReadFile(filePath string) {

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Printf("failed reading config file %s: %v", filePath, err)
		return
	}

	var cfg Configuration
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("failed parsing config file %s: %v\n", filePath, err)
		return
	}

	if len(cfg.FactorioUser) > 0 {
		Loaded.FactorioUser = cfg.FactorioUser
	}

	if len(cfg.FactorioToken) > 0 {
		Loaded.FactorioToken = cfg.FactorioToken
	}

	if len(cfg.WebhookChat) > 0 {
		Loaded.WebhookChat = cfg.WebhookChat
	}

	if len(cfg.WebhookStatus) > 0 {
		Loaded.WebhookStatus = cfg.WebhookStatus
	}

	for _, mod := range cfg.Mods {
		Loaded.Mods = append(Loaded.Mods, mod)
	}

}

func Clean() {

	Loaded.FactorioUser = strings.TrimSpace(Loaded.FactorioUser)

	Loaded.FactorioToken = strings.TrimSpace(Loaded.FactorioToken)

	slices.Sort(Loaded.Mods)
	Loaded.Mods = slices.Compact(Loaded.Mods)
	for i := range Loaded.Mods {
		Loaded.Mods[i] = strings.TrimSpace(Loaded.Mods[i])
	}

}

func Validate() {

	if len(Loaded.FactorioUser) < 1 {
		log.Fatalln(fmt.Sprintf("config missing factorio.com username"))
	}

	if len(Loaded.FactorioToken) < 1 {
		log.Fatalln(fmt.Sprintf("config missing factorio.com token"))
	}

}
