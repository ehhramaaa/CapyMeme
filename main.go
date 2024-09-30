package main

import (
	"CapybaraMeme/core"
	"CapybaraMeme/tools"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

func main() {

	// add driver for support yaml content
	config.AddDriver(yaml.Driver)

	err := config.LoadFiles("configs/config.yml")
	if err != nil {
		panic(err)
	}

	tools.PrintLogo()

	core.LaunchBot()
}
