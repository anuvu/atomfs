package mount

import (
	"os/exec"

	"github.com/anuvu/atomfs/types"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

var MountTypes map[types.AtomType]func(string, string) error

func mountTar(source string, dest string) error {
	cmd := exec.Command("archivemount", source, dest)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Errorf("error mounting %s (%s): %s", source, err, string(output))
	}

	return nil
}

func mountSquashfs(source string, dest string) error {
	err := unix.Mount(source, dest, "squashfs", 0, "")
	return errors.Wrapf(err, "couldn't mount %s to %s", source, dest)
}

func init() {
	MountTypes = map[types.AtomType]func(string, string) error{}
	MountTypes["tar"] = mountTar
	MountTypes["squashfs"] = mountSquashfs
}
