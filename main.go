package main

import (
	"fmt"
	"os"

	"github.com/dylansawicki15/gator/internal/config"
)

func main() {
	configFilePath := config.GetConfigFilePath()

	raw, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return
	}
	fmt.Printf(string(raw) + "\n")

	fileContents, err := config.Read(configFilePath)
	if err != nil {
		fmt.Printf("Error parsing config file: %v\n", err)
		return
	}
	err = config.SetUser(configFilePath, fileContents, "dylan")
	if err != nil {
		fmt.Printf("Error writing config file: %v\n", err)
		return
	}

}
