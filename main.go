package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/dylansawicki15/gator/internal/config"
)

type state struct {
	config *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.handlers[cmd.name]
	if !exists {
		return fmt.Errorf("command not found: %s", cmd.name)
	}

	return handler(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	if c.handlers == nil {
		c.handlers = make(map[string]func(*state, command) error)
	}
	c.handlers[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("a username is required")
	}

	username := cmd.arguments[0]
	configPath := config.GetConfigFilePath()

	if err := config.SetUser(configPath, s.config, username); err != nil {
		return err
	}

	fmt.Printf("Current user has been set to %s\n", username)
	return nil
}

func main() {
	configFilePath := config.GetConfigFilePath()

	if len(os.Args) < 2 {
		fmt.Println("not enough arguments provided")
		os.Exit(1)
	}

	cfg, err := config.Read(configFilePath)
	if err != nil {
		fmt.Printf("error reading config file: %v\n", err)
		os.Exit(1)
	}

	appState := state{config: cfg}
	cmds := commands{handlers: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)

	cmd := command{
		name:      os.Args[1],
		arguments: os.Args[2:],
	}

	if err := cmds.run(&appState, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
