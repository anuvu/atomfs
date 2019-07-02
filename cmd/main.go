package main

import (
	"fmt"
	"os"

	"github.com/anuvu/atomfs/types"
	"github.com/urfave/cli"
)

var (
	version = ""
	debug   = false
)

func getAtomfsConfig(ctx *cli.Context) (types.Config, error) {
	return types.NewConfig(ctx.GlobalString("base-dir"))
}

func main() {
	app := cli.NewApp()
	app.Name = "atomfs"
	app.Usage = "atomfs manages container filesystems"
	app.Version = version
	app.Commands = []cli.Command{
		slurpOCICmd,
		mountCmd,
		umountCmd,
		fsckCmd,
		gcCmd,
		initCmd,
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "base-dir",
			Usage: "the base atomfs dir for managing data",
			Value: "/var/lib/atomfs",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "print stack traces on exceptions",
		},
	}

	app.Before = func(ctx *cli.Context) error {
		debug = ctx.Bool("debug")
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		format := "error: %v\n"
		if debug {
			format = "error: %+v\n"
		}
		fmt.Fprintf(os.Stderr, format, err)
		os.Exit(1)
	}
}
