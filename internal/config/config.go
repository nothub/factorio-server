package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"
)

var FilePath = "config.yaml"
var Loaded Configuration

type Configuration struct {
	WorkDir       string   `yaml:"workdir"`
	FactorioUser  string   `yaml:"factorio-user"`
	FactorioToken string   `yaml:"factorio-token"`
	Mods          []string `yaml:"mods"`
}

func Load() {

	SetDefaults()
	ReadFile(FilePath)
	ReadFlags()
	Clean()
	Validate()

}

func SetDefaults() {
	Loaded.WorkDir = "server"
}

func ReadFile(filePath string) {

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("failed reading config file %s: %v", filePath, err)
		return
	}

	var cfg Configuration
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Printf("failed parsing config file %s: %v", filePath, err)
		return
	}

	if len(cfg.FactorioUser) > 0 {
		Loaded.FactorioUser = cfg.FactorioUser
	}

	if len(cfg.FactorioToken) > 0 {
		Loaded.FactorioToken = cfg.FactorioToken
	}

	for _, mod := range cfg.Mods {
		Loaded.Mods = append(Loaded.Mods, mod)
	}

	if len(cfg.WorkDir) > 0 {
		Loaded.WorkDir = cfg.WorkDir
	}
}

func ReadFlags() {

	workDirP := flag.String("workdir", "", "server base dir and process pwd")
	factorioUserP := flag.String("factorio-user", "", "factorio.com username")
	factorioTokenP := flag.String("factorio-token", "", "factorio.com token")
	modsP := flag.String("mods", "", "game mods")

	flag.Parse()

	if len(*workDirP) > 0 {
		Loaded.WorkDir = *workDirP
	}

	if len(*factorioUserP) > 0 {
		Loaded.FactorioUser = *factorioUserP
	}

	if len(*factorioTokenP) > 0 {
		Loaded.FactorioToken = *factorioTokenP
	}

	if len(*modsP) > 0 {
		for _, mod := range strings.Split(*modsP, ",") {
			if len(mod) > 0 {
				Loaded.Mods = append(Loaded.Mods, mod)
			}
		}
	}

}

func Clean() {

	Loaded.WorkDir = strings.TrimSpace(Loaded.WorkDir)

	Loaded.FactorioUser = strings.TrimSpace(Loaded.FactorioUser)

	Loaded.FactorioToken = strings.TrimSpace(Loaded.FactorioToken)

	slices.Sort(Loaded.Mods)
	Loaded.Mods = slices.Compact(Loaded.Mods)
	for i := range Loaded.Mods {
		Loaded.Mods[i] = strings.TrimSpace(Loaded.Mods[i])
	}

}

func Validate() {

	_, err := filepath.Abs(Loaded.WorkDir)
	if err != nil {
		log.Fatalln(fmt.Sprintf("config workdir invalid: %v", err))
	}

	if len(Loaded.FactorioUser) < 1 {
		log.Fatalln(fmt.Sprintf("config missing factorio.com username"))
	}

	if len(Loaded.FactorioToken) < 1 {
		log.Fatalln(fmt.Sprintf("config missing factorio.com token"))
	}

}
