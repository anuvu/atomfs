package main

import (
	"github.com/anuvu/atomfs"
	"github.com/urfave/cli"
)

var gcCmd = cli.Command{
	Name:   "gc",
	Usage:  "does a garbage collection on an atomfs",
	Action: doGC,
}

func doGC(ctx *cli.Context) error {
	fs, err := atomfs.New(getAtomfsConfig(ctx))
	if err != nil {
		return err
	}
	defer fs.Close()
	return fs.GC()
}
