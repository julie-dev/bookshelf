package main

import (
	"flag"
	"github.com/spf13/pflag"
)

const (
	ApplicationVersion = "v1"
)

type Opts struct {
	Version      		string ``
	ServerURL  	 		string `env:"HTTP_SERVER_URL,required"`
	LogLevel     		string `env:"HTTP_SERVER_LOGLEVEL,default=debug"`
	LogFile      		string `env:"HTTP_SERVER_LOGFILE"`
}

func installFlags(flags *pflag.FlagSet, c *Opts) {
	flags.StringVar(&c.ServerURL, "url", c.ServerURL, "set the atto server ip:port")
	flags.StringVar(&c.LogLevel, "log-level", c.LogLevel, "set the log output level e.g)\"info\",\"debug\"")
	flags.StringVar(&c.LogFile, "log-file", c.LogFile, "set the log file path")

	flagset := flag.NewFlagSet("mariadb", flag.PanicOnError)
	initMariaDBFlags(flagset, c)
	flagset.VisitAll(func(f *flag.Flag) {
		f.Name = "mariadb." + f.Name
		flags.AddGoFlag(f)
	})

	c.Version = ApplicationVersion
}

func initMariaDBFlags(flagset *flag.FlagSet, c *Opts) {
	if flagset == nil {
		flagset = flag.CommandLine
	}

	//flagset.StringVar(&c.MariaDBURL, "url", c.MariaDBURL, "set MariaDB's ip:port")
}
