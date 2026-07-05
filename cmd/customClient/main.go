package main

import (
	"httpLoadTester/internal/config"
	"httpLoadTester/internal/loadtest"
)

func main() {
	cfg, tr := config.Load()
	loadtest.Run(cfg, tr)
}
