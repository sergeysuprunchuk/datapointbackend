package main

import (
	"datapointbackend/config"
	"datapointbackend/internal/app"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err.Error())
	}

	if err = app.Run(cfg); err != nil {
		panic(err.Error())
	}
}
