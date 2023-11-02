package main

import "flag"

var serverDir string

func init() {
	serverDirP := flag.String("server-dir", "server", "Server base dir and process pwd")
	flag.Parse()
	serverDir = *serverDirP
}
