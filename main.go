package main

import (
	"os"

	"github.com/RofaBR/usdt-monitoring-svc/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
