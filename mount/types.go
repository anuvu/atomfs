package mount

import (
	"os/exec"

	"github.com/freddierice/go-losetup"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

var MountTypes map[string]func(string, string) error

func mountTar(source string, dest string) error {
	cmd := exec.Command("archivemount", source, dest)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.Errorf("error mounting %s (%s): %s", source, err, string(output))
	}

	return nil
}

func mountSquashfs(source string, dest string) error {
	dev, err := losetup.Attach(source, 0, true)
	if err != nil {
		return err
	}
	defer dev.Detach()

	err = unix.Mount(dev.Path(), dest, "squashfs", unix.MS_RDONLY, "")
	return errors.Wrapf(err, "couldn't mount %s to %s", source, dest)
}

func init() {
	MountTypes = map[string]func(string, string) error{}
	MountTypes["tar"] = mountTar
	MountTypes["squashfs"] = mountSquashfs
}
