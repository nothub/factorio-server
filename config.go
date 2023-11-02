package main

import (
	"flag"
	"log"
	"path/filepath"
)

var serverDir string
var factorioUser string
var factorioToken string

func init() {
	serverDirP := flag.String("server-dir", "server", "Server base dir and process pwd")
	factorioUserP := flag.String("factorio-user", "", "factorio.com username")
	factorioTokenP := flag.String("factorio-token", "", "factorio.com token")

	flag.Parse()

	p, err := filepath.Abs(*serverDirP)
	if err != nil {
		log.Fatalln(err)
	}
	serverDir = p

	factorioUser = *factorioUserP
	factorioToken = *factorioTokenP
}
