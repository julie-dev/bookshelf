package main

import (
	"bookshelf/config"
	"bookshelf/controller"
	"bookshelf/database"
	"bookshelf/router"
	"context"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/joeshaw/envdecode"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

// @title Bookshelf
// @version 1.0
// @description System for register and manage book list

// @host localhost:8080
// @BasePath /api/v1/books

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

	repository, err := database.NewRepository(&config.Database)
	if err != nil {
		return err
	}

	service := controller.NewBookshelfService(&config, repository)

	e := router.New(service, c.APIVersion)
	e.Server.Addr = c.ServerURL

	if c.LogLevel == "debug" {
		e.Debug = true
	}

	e.Logger.Fatal(gracehttp.Serve(e.Server))

	return nil
}
