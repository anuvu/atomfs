package atomfs

import (
	"context"
	"path"

	"github.com/anuvu/atomfs/types"
	"github.com/opencontainers/umoci"
	stackeroci "github.com/anuvu/stacker/oci"
)

type MountOCIOpts struct {
	OCIDir       string
	MetadataPath string
	Tag          string
	Target       string
	Writable     bool
}

func (c MountOCIOpts) AtomsPath(parts ...string) string {
	atoms := path.Join(c.OCIDir, "blobs", "sha256")
	return path.Join(append([]string{atoms}, parts...)...)
}

func (c MountOCIOpts) MountedAtomsPath(parts ...string) string {
	mounts := path.Join(c.MetadataPath, "mounts")
	return path.Join(append([]string{mounts}, parts...)...)
}

func (c MountOCIOpts) OverlayDirsPath(parts ...string) string {
	overlayDirs := path.Join(c.MetadataPath, "overlay-dirs")
	return path.Join(append([]string{overlayDirs}, parts...)...)
}

func BuildMoleculeFromOCI(opts MountOCIOpts) (types.Molecule, error) {
	oci, err := umoci.OpenLayout(opts.OCIDir)
	if err != nil {
		return types.Molecule{}, err
	}
	defer oci.Close()

	man, err := stackeroci.LookupManifest(oci, opts.Tag)
	if err != nil {
		return types.Molecule{}, err
	}

	atoms := []types.Atom{}
	for _, l := range man.Layers {
		layer, err := oci.FromDescriptor(context.Background(), l)
		if err != nil {
			return types.Molecule{}, err
		}
		defer layer.Close()

		atom := types.Atom{Hash: l.Digest.Encoded(), Type: types.SquashfsAtom}
		atoms = append(atoms, atom)
	}

	// The OCI spec says that the first layer should be the bottom most
	// layer. In overlay it's the top most layer. Since the atomfs codebase
	// is mostly a wrapper around overlayfs, let's keep things in our db in
	// the same order that overlay expects them, i.e. the first layer is
	// the top most. That means we need to reverse the order in which the
	// atoms were inserted, because they were backwards.
	//
	// It's also terrible that golang doesn't have a reverse function, but
	// that's a discussion for a different block comment.
	for i := len(atoms)/2 - 1; i >= 0; i-- {
		opp := len(atoms) - 1 - i
		atoms[i], atoms[opp] = atoms[opp], atoms[i]
	}

	mol := types.NewMolecule(opts)
	mol.Name = opts.Tag
	mol.Atoms = atoms

	return mol, nil
}

func UnmountOCI(opts MountOCIOpts) error {
	return types.Umount(opts, opts.Target)
}
