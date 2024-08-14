package main

import (
	"log"
	"os"

	"github.com/RofaBR/usdt-monitoring-svc/internal/cli"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
