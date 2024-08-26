package main

import (
	"os"

	"github.com/RofaBR/usdt-monitoring-svc/internal/cli"
)

func main() {
	/*err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}*/
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
