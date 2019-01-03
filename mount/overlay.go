package mount

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/anuvu/atomfs/types"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

type Overlay struct {
	config types.Config
	atoms  []types.Atom
}

func NewOverlay(config types.Config, atoms []types.Atom) (*Overlay, error) {
	return &Overlay{config: config, atoms: atoms}, nil
}

func (o *Overlay) Mount(dest string) error {
	// The kernel unfortunately doesn't support mntopts > 4096 characters,
	// so let's figure out if we've got too many atoms here:
	//     len("lowerdir=") + len(o.atoms) * (len(config.Path) + len("/atoms/") + 64 + 1)
	// * 64 is the len(sha256sum) + 1 for the : separator
	charCount := len("lowerdir=") + len(o.atoms)*(len(o.config.Path)+len("/atoms/")+64+1)
	if charCount > 4096 {
		return fmt.Errorf("too many lower dirs; must have fewer than 4096 chars")
	}

	dirs := []string{}
	// first, mount everything
	for _, a := range o.atoms {
		target := o.config.MountedAtomsPath(a.Hash)
		dirs = append(dirs, target)
		_, err := os.Stat(target)
		if err == nil {
			continue
		}

		if !os.IsNotExist(err) {
			return err
		}

		if err := os.MkdirAll(target, 755); err != nil {
			return err
		}

		mounter, ok := MountTypes[a.Type]
		if !ok {
			return errors.Errorf("don't know how to mount %s of type %s", a.Name, a.Type)
		}

		if err := mounter(o.config.AtomsPath(a.Hash), target); err != nil {
			return errors.Wrapf(err, "couldn't mount")
		}
	}

	// overlay doesn't work with one lowerdir. so we do a hack here: we
	// just create an empty directory called "workaround" in the mounts
	// directory, and add that to the dir list if it's of length one.
	if len(dirs) == 1 {
		workaround := o.config.MountedAtomsPath("workaround")
		if err := os.MkdirAll(workaround, 755); err != nil {
			return errors.Wrapf(err, "couldn't make workaround dir")
		}

		dirs = append(dirs, workaround)
	}

	// now, do the actual overlay mount
	err := unix.Mount("overlay", dest, "overlay", 0, "lowerdir="+strings.Join(dirs, ":"))
	return errors.Wrapf(err, "couldn't do overlay mount")
}

func getOverlayDirs(m Mount) []string {
	for _, opt := range m.Opts {
		if !strings.HasPrefix(opt, "lowerdir=") {
			continue
		}

		return strings.Split(strings.TrimPrefix(opt, "lowerdir="), ":")
	}

	return []string{}
}

func Umount(dest string) error {
	mounts, err := ParseMounts()
	if err != nil {
		return err
	}

	underlyingAtoms := []string{}
	for _, m := range mounts {
		fmt.Println("checking", m.FSType, "mount", m.Target)
		if m.Target != dest || m.FSType != "overlay" {
			continue
		}

		underlyingAtoms = getOverlayDirs(m)
	}

	if len(underlyingAtoms) == 0 {
		return errors.Errorf("%s is not an atomfs mountpoint", dest)
	}

	if err := unix.Unmount(dest, 0); err != nil {
		return err
	}

	// now, "refcount" the remaining atoms and see if any of ours are
	// unused
	usedAtoms := map[string]bool{}

	mounts, err = ParseMounts()
	if err != nil {
		return err
	}

	for _, m := range mounts {
		if m.FSType != "overlay" {
			continue
		}

		dirs := getOverlayDirs(m)
		for _, d := range dirs {
			usedAtoms[d] = true
		}
	}

	// If any of the atoms underlying the target mountpoint are now unused,
	// let's unmount them too.
	for _, a := range underlyingAtoms {
		_, used := usedAtoms[a]
		if used {
			continue
		}

		// the workaround dir isn't really a mountpoint, so don't unmount it
		if path.Base(a) == "workaround" {
			continue
		}

		if err := unix.Unmount(a, 0); err != nil {
			return err
		}
	}

	return nil
}

type Mount struct {
	Source string
	Target string
	FSType string
	Opts   []string
}

func ParseMounts() ([]Mount, error) {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	mounts := []Mount{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		mount := Mount{}
		mount.Target = fields[4]

		for i := 5; i < len(fields); i++ {
			if fields[i] != "-" {
				continue
			}

			mount.FSType = fields[i+1]
			mount.Source = fields[i+2]
			mount.Opts = strings.Split(fields[i+3], ",")
		}

		mounts = append(mounts, mount)
	}

	return mounts, nil
}
