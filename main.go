package main

import (
	"fmt"

	"github.com/nordluma/gator/internal/config"
)

func main() {
	cfg := config.Read()
	cfg.SetUser("lane")
	cfg = config.Read()

	fmt.Println(cfg)
}
