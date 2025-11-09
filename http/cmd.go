package main

import (
	"fmt"
	"playgoround/http/internal/config"
)

func main() {
	cfg := config.Load()
	fmt.Printf("load config success: %+v", cfg)

	app := New(cfg)
	app.Start()
}
