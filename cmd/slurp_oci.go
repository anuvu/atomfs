package main

import (
	"github.com/anuvu/atomfs"
	"github.com/urfave/cli"
)

var slurpOCICmd = cli.Command{
	Name:   "slurp-oci",
	Usage:  "import an OCI image to atomfs",
	Action: doSlurpOCI,
	ArgsUsage: `<oci-dir>

Import the OCI directory into atomfs. Note that this copies the OCI blobs to
the atomfs directory.
`,
}

func doSlurpOCI(ctx *cli.Context) error {
	fs, err := atomfs.New(getAtomfsConfig(ctx))
	if err != nil {
		return err
	}
	defer fs.Close()
	return fs.SlurpOCI(ctx.Args().Get(0))
}
