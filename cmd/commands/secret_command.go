package commands

import (
	"flag"
	"fmt"

	"httpkeeper/token"
)

// NewSecretCommand creates command to generate secret
func NewSecretCommand() *SecretCommand {
	sc := &SecretCommand{
		fs: flag.NewFlagSet("secret", flag.ExitOnError),
	}

	sc.fs.IntVar(&sc.length, "len", 128, "lenght of the secret")
	return sc
}

// SecretCommand contains flags related to secret command
type SecretCommand struct {
	fs     *flag.FlagSet
	length int
}

// Init initializes flagset with args
func (sc *SecretCommand) Init(args []string) error {
	return sc.fs.Parse(args)
}

// Name returns command name
func (sc *SecretCommand) Name() string {
	return "secret"
}

// Run runs secret generating
func (sc *SecretCommand) Run() error {
	secret, err := token.GenerateSecret(sc.length)
	if err != nil {
		return err
	}
	fmt.Println(secret)
	return nil
}
