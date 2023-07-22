package commands

import (
	"flag"
	"fmt"
	"os"
	"time"

	"httpkeeper/config"
	"httpkeeper/token"
)

const expiresAtLayout = "2006-01-02 15:04:05"

// NewTokenCommand creates command to generate token
func NewTokenCommand() *TokenCommand {
	tg := &TokenCommand{
		fs: flag.NewFlagSet("token", flag.ExitOnError),
	}

	tg.fs.StringVar(&tg.client, "client", "client", "client name")
	tg.fs.StringVar(&tg.config, "config", "./config.json", "path to configuration file")
	tg.fs.StringVar(&tg.expiresAt, "expiresAt",
		time.Now().AddDate(0, 0, 1).Format(expiresAtLayout),
		"time thats token valid until, by default it is set 1 day from now (format: yyyy-mm-dd HH:MM:SS)",
	)

	return tg
}

// TokenCommand contains flags related to token command
type TokenCommand struct {
	fs        *flag.FlagSet
	client    string
	config    string
	expiresAt string
}

// Init initializes flagset with args
func (tg *TokenCommand) Init(args []string) error {
	return tg.fs.Parse(args)
}

// Name returns command name
func (tg *TokenCommand) Name() string {
	return "token"
}

// Run runs token generating
func (tg *TokenCommand) Run() error {
	expiresAt, err := time.Parse(expiresAtLayout, tg.expiresAt)
	if err != nil {
		return err
	}

	conf, err := config.Load(tg.config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot load config: %v\n", err)
		os.Exit(1)
	}

	clientToken, err := token.Generate(conf.Server.Secret, conf.Server.Name, tg.client, expiresAt)
	if err != nil {
		return err
	}
	fmt.Println(clientToken)
	return nil
}
