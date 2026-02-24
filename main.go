package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dylansawicki15/gator/internal/config"
	"github.com/dylansawicki15/gator/internal/database"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
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

	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		fmt.Println("user does not exist")
		os.Exit(1)
	}

	configPath := config.GetConfigFilePath()
	if err := config.SetUser(configPath, s.cfg, username); err != nil {
		return err
	}

	fmt.Printf("Current user has been set to %s\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("a name is required")
	}

	name := cmd.arguments[0]
	now := time.Now()

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
	})
	if err != nil {
		fmt.Println("user already exists")
		os.Exit(1)
	}

	configPath := config.GetConfigFilePath()
	if err := config.SetUser(configPath, s.cfg, name); err != nil {
		return err
	}

	fmt.Printf("User created: %+v\n", user)
	return nil
}

func handlerReset(s *state, _ command) error {
	if err := s.db.DeleteUsers(context.Background()); err != nil {
		fmt.Printf("error resetting database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("database reset successfully")
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

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		fmt.Printf("error opening database: %v\n", err)
		os.Exit(1)
	}

	dbQueries := database.New(db)

	appState := state{
		db:  dbQueries,
		cfg: cfg,
	}
	cmds := commands{handlers: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)

	cmd := command{
		name:      os.Args[1],
		arguments: os.Args[2:],
	}

	if err := cmds.run(&appState, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
