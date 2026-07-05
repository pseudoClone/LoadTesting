package main

import (
	"httpLoadTester/internal/config"
	"httpLoadTester/internal/loadtest"
)

func main() {
	cfg := config.Load()
	loadtest.Run(cfg)
}
