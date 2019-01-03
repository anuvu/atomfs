package main

import (
	"github.com/anuvu/atomfs"
	"github.com/urfave/cli"
)

var umountCmd = cli.Command{
	Name:    "umount",
	Aliases: []string{"unmount"},
	Usage:   "mounts an atomfs molecule at a location",
	Action:  doUmount,
	ArgsUsage: `<mountpoint>

unmounts the specified mountpoint, cleaning up intermediate atom mounts as
applicable.
`,
}

func doUmount(ctx *cli.Context) error {
	fs, err := atomfs.New(getAtomfsConfig(ctx))
	if err != nil {
		return err
	}
	defer fs.Close()
	return fs.Umount(ctx.Args().Get(0))
}
