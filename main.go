package main

import (
	"fmt"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"

	"capybara-meme/core"
)

func main() {
	fmt.Println(`
  /$$$$$$                                      /$$      /$$                                  
 /$$__  $$                                    | $$$    /$$$                                  
| $$  \__/  /$$$$$$   /$$$$$$  /$$   /$$      | $$$$  /$$$$  /$$$$$$  /$$$$$$/$$$$   /$$$$$$ 
| $$       |____  $$ /$$__  $$| $$  | $$      | $$ $$/$$ $$ /$$__  $$| $$_  $$_  $$ /$$__  $$
| $$        /$$$$$$$| $$  \ $$| $$  | $$      | $$  $$$| $$| $$$$$$$$| $$ \ $$ \ $$| $$$$$$$$
| $$    $$ /$$__  $$| $$  | $$| $$  | $$      | $$\  $ | $$| $$_____/| $$ | $$ | $$| $$_____/
|  $$$$$$/|  $$$$$$$| $$$$$$$/|  $$$$$$$      | $$ \/  | $$|  $$$$$$$| $$ | $$ | $$|  $$$$$$$
 \______/  \_______/| $$____/  \____  $$      |__/     |__/ \_______/|__/ |__/ |__/ \_______/
                    | $$       /$$  | $$                                                     
                    | $$      |  $$$$$$/                                                     
                    |__/       \______/                                                      
`)
	fmt.Println(`ρσωєяє∂ ву : нσℓу¢αη`)

	// add driver for support yaml content
	config.AddDriver(yaml.Driver)

	err := config.LoadFiles("config.yml")
	if err != nil {
		panic(err)
	}

	core.ProcessBot(config.Default())
}
