package main

import (
	"fmt"
	"msu-logging-backend/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
}
