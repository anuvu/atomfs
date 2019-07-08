package main

import (
	"io"
	"os"

	"github.com/anuvu/atomfs"
	"github.com/urfave/cli"
)

var dumpDBCmd = cli.Command{
	Name:   "dump-db",
	Usage:  "initializes an empty atomfs in the target directory",
	Action: doDumpDB,
	Hidden: true,
}

func doDumpDB(ctx *cli.Context) error {
	config, err := getAtomfsConfig(ctx)
	if err != nil {
		return err
	}

	fs, err := atomfs.New(config)
	if err != nil {
		return err
	}
	defer fs.Close()

	_, err = io.Copy(os.Stdout, fs.DumpDB())
	return err
}
