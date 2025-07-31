package config

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

var FactorioUser string
var FactorioToken string
var Mods []string

func Load() {

	ReadFile("config.yaml")
	ReadFlags()
	Validate()

}

func ReadFile(s string) {

	// TODO

}

func ReadFlags() {

	factorioUserP := flag.String("factorio-user", "", "factorio.com username")
	factorioTokenP := flag.String("factorio-token", "", "factorio.com token")
	factorioModsP := flag.String("mods", "", "game mods")

	flag.Parse()

	if len(strings.TrimSpace(*factorioUserP)) > 0 {
		FactorioUser = *factorioUserP
	}

	if len(strings.TrimSpace(*factorioTokenP)) > 0 {
		FactorioToken = *factorioTokenP
	}

	{
		var factorioMods []string
		for _, mod := range strings.Split(*factorioModsP, ",") {
			if len(strings.TrimSpace(mod)) > 0 {
				factorioMods = append(factorioMods, mod)
			}
		}
		Mods = factorioMods
	}

}

func Validate() {

	if len(strings.TrimSpace(FactorioUser)) < 1 {
		log.Fatalln(fmt.Sprintf("config missing factorio.com username"))
	}

	if len(strings.TrimSpace(FactorioToken)) < 1 {
		log.Fatalln(fmt.Sprintf("config missing factorio.com token"))
	}

}
