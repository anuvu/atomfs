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

type Config interface {
	AtomsPath(parts ...string) string
	MountedAtomsPath(parts ...string) string
	OverlayDirsPath(parts ...string) string
}

type DBBasedConfig struct {
	Path string
}

func NewDBBasedConfig(path string) (DBBasedConfig, error) {
	config := DBBasedConfig{path}

	if err := os.MkdirAll(config.AtomsPath(), 0755); err != nil {
		return DBBasedConfig{}, err
	}
	if err := os.MkdirAll(config.MountedAtomsPath(), 0755); err != nil {
		return DBBasedConfig{}, err
	}
	if err := os.MkdirAll(config.OverlayDirsPath(), 0755); err != nil {
		return DBBasedConfig{}, err
	}

	return config, nil
}

func (c DBBasedConfig) RelativePath(parts ...string) string {
	return path.Join(append([]string{c.Path}, parts...)...)
}

func (c DBBasedConfig) AtomsPath(parts ...string) string {
	atoms := c.RelativePath("atoms")
	return path.Join(append([]string{atoms}, parts...)...)
}

func (c DBBasedConfig) MountedAtomsPath(parts ...string) string {
	mounts := c.RelativePath("mounts")
	return path.Join(append([]string{mounts}, parts...)...)
}

func (c DBBasedConfig) OverlayDirsPath(parts ...string) string {
	overlayDirs := c.RelativePath("overlay-dirs")
	return path.Join(append([]string{overlayDirs}, parts...)...)
}

type AtomType string

const (
	TarAtom      AtomType = "tar"
	SquashfsAtom AtomType = "squashfs"
)
