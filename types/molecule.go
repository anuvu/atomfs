package types

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/anuvu/atomfs/mount"
	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

type Molecule struct {
	ID   int64
	Name string
	// Atoms is the list of atoms in this Molecule. The first element in
	// this list is the top most layer in the overlayfs.
	Atoms []Atom

	config Config
}

func NewMolecule(config Config) Molecule {
	return Molecule{config: config}
}

// MountUnderlyingAtoms mounts all the underlying atoms at
// config.MountedAtomsPath().
func (m Molecule) MountUnderlyingAtoms() error {
	mounts, err := mount.ParseMounts("/proc/self/mountinfo")
	if err != nil {
		return errors.Wrapf(err, "couldn't parse mounts")
	}

	dirs := []string{}
	for _, a := range m.Atoms {
		target := m.config.MountedAtomsPath(a.Hash)
		dirs = append(dirs, target)

		if mounts.IsMountpoint(target) {
			continue
		}

		if err := os.MkdirAll(target, 755); err != nil {
			return err
		}

		mounter, ok := mount.MountTypes[string(a.Type)]
		if !ok {
			return errors.Errorf("don't know how to mount %s of type %s", a.Name, a.Type)
		}

		if err := mounter(m.config.AtomsPath(a.Hash), target); err != nil {
			return errors.Wrapf(err, "couldn't mount")
		}
	}

	return nil
}

func sha256string(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// OverlayArgs returns all of the mount options to pass to the kernel to
// actually mount this molecule.
func (m Molecule) OverlayArgs(dest string, writable bool) (string, error) {
	dirs := []string{}
	for _, a := range m.Atoms {
		target := m.config.MountedAtomsPath(a.Hash)
		dirs = append(dirs, target)
	}

	// overlay doesn't work with one lowerdir. so we do a hack here: we
	// just create an empty directory called "workaround" in the mounts
	// directory, and add that to the dir list if it's of length one.
	if len(dirs) == 1 {
		workaround := m.config.MountedAtomsPath("workaround")
		if err := os.MkdirAll(workaround, 755); err != nil {
			return "", errors.Wrapf(err, "couldn't make workaround dir")
		}

		dirs = append(dirs, workaround)
	}

	// Note that in overlayfs, the first thing is the top most layer in the
	// overlay.
	mntOpts := "index=off,lowerdir=" + strings.Join(dirs, ":")
	if writable {
		// In order to make it so that we can Unmount() without saving
		// any state, we construct special names for the workdir and
		// upperdir:
		//   sha256(dest)/{upperdir|workdir}
		// Note that if this already exists, we don't want to re-use it
		// (and indeed we can't, overlay will fail the mount); this
		// means that there can only ever be one atomfs mount at a
		// particular location. That doesn't seem too big a deal, though.
		upperDir := m.config.OverlayDirsPath(sha256string(dest), "upperdir")
		workDir := m.config.OverlayDirsPath(sha256string(dest), "workdir")

		if err := os.MkdirAll(upperDir, 0755); err != nil {
			return "", err
		}
		if err := os.MkdirAll(workDir, 0755); err != nil {
			return "", err
		}

		mntOpts += fmt.Sprintf(",upperdir=%s,workdir=%s", upperDir, workDir)
	}

	return mntOpts, nil
}

func (m Molecule) Mount(dest string, writable bool) error {
	mntOpts, err := m.OverlayArgs(dest, writable)
	if err != nil {
		return err
	}

	// The kernel doesn't allow mount options longer than 4096 chars, so
	// let's give a nicer error than -EINVAL here.
	if len(mntOpts) > 4096 {
		return errors.Errorf("too many lower dirs; must have fewer than 4096 chars")
	}

	err = m.MountUnderlyingAtoms()
	if err != nil {
		return err
	}

	// now, do the actual overlay mount
	err = unix.Mount("overlay", dest, "overlay", 0, mntOpts)
	return errors.Wrapf(err, "couldn't do overlay mount to %s, opts: %s", dest, mntOpts)
}

func Umount(config Config, dest string) error {
	mounts, err := mount.ParseMounts("/proc/self/mountinfo")
	if err != nil {
		return err
	}

	underlyingAtoms := []string{}
	for _, m := range mounts {
		if m.Target != dest || m.FSType != "overlay" {
			continue
		}

		underlyingAtoms = mount.GetOverlayDirs(m)
		break
	}

	if len(underlyingAtoms) == 0 {
		return errors.Errorf("%s is not an atomfs mountpoint", dest)
	}

	if err := unix.Unmount(dest, 0); err != nil {
		return err
	}

	// If this was writable, we should clean up the work/upperdir.
	err = os.RemoveAll(config.OverlayDirsPath(sha256string(dest)))
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// now, "refcount" the remaining atoms and see if any of ours are
	// unused
	usedAtoms := map[string]int{}

	mounts, err = mount.ParseMounts("/proc/self/mountinfo")
	if err != nil {
		return err
	}

	for _, m := range mounts {
		if m.FSType != "overlay" {
			continue
		}

		dirs := mount.GetOverlayDirs(m)
		for _, d := range dirs {
			usedAtoms[d]++
		}
	}

	// If any of the atoms underlying the target mountpoint are now unused,
	// let's unmount them too.
	for _, a := range underlyingAtoms {
		count, used := usedAtoms[a]
		if used && count > 1 {
			continue
		}
		/* TODO: some kind of logging
		if !used {
			log.Warnf("unused atom %s was part of this molecule?")
			continue
		}
		*/

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
