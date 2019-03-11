package main

import (
	"github.com/anuvu/atomfs"
	"github.com/urfave/cli"
)

var gcCmd = cli.Command{
	Name:  "gc",
	Usage: "does a garbage collection on an atomfs",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "dry-run",
			Usage: "do a dry run of a GC, without actually deleting anything",
		},
	},
	Action: doGC,
}

func doGC(ctx *cli.Context) error {
	config, err := getAtomfsConfig(ctx)
	if err != nil {
		return err
	}

	fs, err := atomfs.New(config)
	if err != nil {
		return err
	}
	defer fs.Close()
	return fs.GC(ctx.Bool("dry-run"))
}
