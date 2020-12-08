package main

import (
	"bookshelf/config"
	"bookshelf/controller"
	"bookshelf/router"
	"context"
	"github.com/joeshaw/envdecode"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	buildVersion = "1.0"
	buildTime    = "20201225"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		cancel()
	}()

	fileName := filepath.Base(os.Args[0])

	var opts Opts
	envdecode.MustDecode(&opts)
	opts.Version = strings.Join([]string{fileName, buildVersion, buildTime}, "-")

	cmd := NewCommand(ctx, fileName, opts)
	cmd.Execute()

	return
}

func NewCommand(ctx context.Context, name string, c Opts) *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: name + " provides a authentication interface for kubernetes cluster nodes.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(ctx, c)
		},
	}

	installFlags(cmd.Flags(), &c)
	return cmd
}

func Run(ctx context.Context, c Opts) error {

	var config config.Config
	envdecode.MustDecode(&config)
	service := controller.NewBookshelfService(&config)

	e := router.New(service)

	if c.LogLevel == "debug" {
		e.Debug = true
	}

	e.Logger.Fatal(e.Start(c.ServerURL))

	return nil
}
