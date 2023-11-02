package main

import (
	"flag"
	"log"
)

func main() {
	log.SetFlags(0)
	log.Println("The factory must grow!")
	run(flag.Args())
}
