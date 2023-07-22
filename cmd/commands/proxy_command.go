package commands

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"httpkeeper/config"
	"httpkeeper/proxy"
)

// NewProxyCommand creates command to run proxy server
func NewProxyCommand() *ProxyCommand {
	pc := &ProxyCommand{
		fs: flag.NewFlagSet("proxy", flag.ExitOnError),
	}

	pc.fs.StringVar(&pc.addr, "addr", "localhost:3000", "addres of the proxy service")
	pc.fs.StringVar(&pc.config, "config", "./config.json", "path to configuration file")
	pc.fs.StringVar(&pc.logFile, "log", "./httpkeeper.log", "path to log file")

	return pc
}

// ProxyCommand contains flags related to proxy command
type ProxyCommand struct {
	fs      *flag.FlagSet
	addr    string
	config  string
	logFile string
}

// Init initializes flagset with args
func (pc *ProxyCommand) Init(args []string) error {
	return pc.fs.Parse(args)
}

// Name returns command name
func (pc *ProxyCommand) Name() string {
	return "proxy"
}

func (pc *ProxyCommand) loadConfig(server *proxy.Proxy) error {
	conf, err := config.Load(pc.config)
	if err != nil {
		return err
	}
	server.SetSecret(conf.Server.Secret)
	server.SetServiceName(conf.Server.Name)
	server.SetServices(conf.Services)
	server.SetInvalidatedTokens(conf.InvalidatedTokens)
	return nil
}

// Run starts the proxy server.
func (pc *ProxyCommand) Run() error {
	conf, err := config.Load(pc.config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot load config: %v\n", err)
		os.Exit(1)
	}

	logFile, err := os.OpenFile(conf.Server.LogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.Printf("starting proxy on %v", conf.Server.Addr)

	server := proxy.New(conf.Server.Addr)
	server.SetSecret(conf.Server.Secret)
	server.SetServiceName(conf.Server.Name)
	server.SetServices(conf.Services)
	server.SetInvalidatedTokens(conf.InvalidatedTokens)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	defer cancel()

	reloadSigChan := make(chan os.Signal)
	signal.Notify(reloadSigChan, syscall.SIGUSR1)

	serverErr := make(chan error)
	go func() {
		serverErr <- server.ListenAndServe()
	}()
	for {
		select {
		case <-reloadSigChan:
			log.Printf("try to reload config...")
			err := pc.loadConfig(server)
			if err != nil {
				log.Printf("loading config failed: %v", err)
			}
			log.Printf("config is reloaded")
		case <-ctx.Done():
			closeCtx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.Server.ShutdownTimeout)*time.Second)
			defer cancel()
			log.Printf("shutting down server, timeout: %vs", conf.Server.ShutdownTimeout)
			if err := server.Shutdown(closeCtx); err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					log.Printf("Shutdown timeout deadline exceeded")
				}
			}
			os.Exit(0)
		case err := <-serverErr:
			log.Printf("server failed to bind on %v, error: %v", conf.Server.Addr, err)
			os.Exit(1)
		}
	}
	return nil
}
