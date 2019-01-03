package mount

import (
	"os/exec"

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
	return unix.Mount(source, dest, "squashfs", 0, "")
}

func init() {
	MountTypes = map[string]func(string, string) error{}
	MountTypes["tar"] = mountTar
	MountTypes["squashfs"] = mountSquashfs
}
