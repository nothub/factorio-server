package main

import (
	"flag"
	"github.com/nothub/factorio-server/internal/config"
	"github.com/nothub/factorio-server/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFlags(0)
	log.Println("The factory must grow!")

	config.Load()

	err := os.MkdirAll("server", 0755)
	if err != nil {
		log.Fatalf("failed creating server dir: %v\n", err)
		return
	}

	err = os.Chdir("server")
	if err != nil {
		log.Fatalf("failed switching to server dir: %v\n", err)
		return
	}

	quit := server.Run(flag.Args())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	for {
		s := <-signals
		switch s {
		case syscall.SIGINT, syscall.SIGTERM:
			log.Printf("Signal %s received, shutting down server...\n", s.String())
			quit()
			os.Exit(0)
		default:
			log.Printf("Signal %s received, handling is not implemented...\n", s.String())
		}
	}
}
