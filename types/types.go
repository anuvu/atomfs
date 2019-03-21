// To avoid circular imports between db/mount/etc. packages, let's define the
// basic atomfs types in their own package that everyone will then import.
package types

import (
	"os"
	"path"
)

type Atom struct {
	ID   int64
	Name string
	Hash string
	Type AtomType
}

type Molecule struct {
	ID   int64
	Name string
	// Atoms is the list of atoms in this Molecule. The first element in
	// this list is the top most layer in the overlayfs.
	Atoms []Atom
}

type Config struct {
	Path string
}

func NewConfig(path string) (Config, error) {
	config := Config{path}

	if err := os.MkdirAll(config.AtomsPath(), 0755); err != nil {
		return Config{}, err
	}
	if err := os.MkdirAll(config.MountedAtomsPath(), 0755); err != nil {
		return Config{}, err
	}
	if err := os.MkdirAll(config.OverlayDirsPath(), 0755); err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c Config) RelativePath(parts ...string) string {
	return path.Join(append([]string{c.Path}, parts...)...)
}

func (c Config) AtomsPath(parts ...string) string {
	atoms := c.RelativePath("atoms")
	return path.Join(append([]string{atoms}, parts...)...)
}

func (c Config) MountedAtomsPath(parts ...string) string {
	mounts := c.RelativePath("mounts")
	return path.Join(append([]string{mounts}, parts...)...)
}

func (c Config) OverlayDirsPath(parts ...string) string {
	overlayDirs := c.RelativePath("overlay-dirs")
	return path.Join(append([]string{overlayDirs}, parts...)...)
}

type AtomType string

const (
	TarAtom      AtomType = "tar"
	SquashfsAtom AtomType = "squashfs"
)
