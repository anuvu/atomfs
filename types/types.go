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
	Type string
}

type Molecule struct {
	ID    int64
	Name  string
	Atoms []Atom
}

type Config struct {
	Path string
}

func (c Config) RelativePath(parts ...string) string {
	return path.Join(append([]string{c.Path}, parts...)...)
}

func (c Config) AtomsPath(parts ...string) string {
	atoms := c.RelativePath("atoms")
	// We explicitly ignore the error here, because we expect it will
	// succeed and it'll fail shortly thereafter if it doesn't. this way we
	// don't have to add a *whole bunch* of error handling code where it's
	// mostly unnecessary.
	os.MkdirAll(atoms, 0755)
	return path.Join(append([]string{atoms}, parts...)...)
}

func (c Config) MountedAtomsPath(parts ...string) string {
	mounts := c.RelativePath("mounts")
	// see note above.
	os.MkdirAll(mounts, 0755)
	return path.Join(append([]string{mounts}, parts...)...)
}
