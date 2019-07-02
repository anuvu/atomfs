package main

import (
	"github.com/anuvu/atomfs"
	"github.com/urfave/cli"
)

var initCmd = cli.Command{
	Name:   "init",
	Usage:  "initializes an empty atomfs in the target directory",
	Action: doInit,
	Hidden: true,
}

func doInit(ctx *cli.Context) error {
	config, err := getAtomfsConfig(ctx)
	if err != nil {
		return err
	}

	fs, err := atomfs.New(config)
	if err != nil {
		return err
	}
	defer fs.Close()
	return nil
}
