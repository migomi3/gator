package main

import (
	"fmt"

	"github.com/migomi3/gator/internal/config"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		fmt.Println(err)
	}

	config.SetUser(&cfg, "Wolf")
	if err != nil {
		fmt.Println(err)
	}

	cfg, err = config.ReadConfig()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(cfg)
}
