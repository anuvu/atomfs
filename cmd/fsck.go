package main

import (
	"fmt"

	"github.com/anuvu/atomfs"
	"github.com/urfave/cli"
)

var fsckCmd = cli.Command{
	Name:   "fsck",
	Usage:  "checks an atomfs filesystem for consistency",
	Action: doFSCK,
}

func doFSCK(ctx *cli.Context) error {
	fs, err := atomfs.New(getAtomfsConfig(ctx))
	if err != nil {
		return err
	}
	defer fs.Close()
	errs, err := fs.FSCK()
	if err != nil {
		return err
	}

	for _, anErr := range errs {
		fmt.Println(anErr)
	}

	if len(errs) > 0 {
		return fmt.Errorf("fsck failed.")
	}

	fmt.Println("fsck ok.")
	return nil
}
