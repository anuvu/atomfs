package main

import (
	"github.com/anuvu/atomfs"
	"github.com/urfave/cli"
)

var mountCmd = cli.Command{
	Name:   "mount",
	Usage:  "mounts an atomfs molecule at a location",
	Action: doMount,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "writable",
			Usage: "set up a writable layer on top",
		},
	},
	ArgsUsage: `<molecule> <mountpoint>

mounts the specified molecule to the specified mountpoint.
`,
}

func doMount(ctx *cli.Context) error {
	config, err := getAtomfsConfig(ctx)
	if err != nil {
		return err
	}

	fs, err := atomfs.New(config)
	if err != nil {
		return err
	}
	defer fs.Close()
	return fs.Mount(ctx.Args().Get(0), ctx.Args().Get(1), ctx.Bool("writable"))
}
