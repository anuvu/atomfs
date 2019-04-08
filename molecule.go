package atomfs

import (
	"context"

	"github.com/anuvu/atomfs/types"
	"github.com/anuvu/stacker"
	"github.com/openSUSE/umoci/oci/casext"
)

func (atomfs *Instance) CreateMolecule(name string, atoms []types.Atom) (types.Molecule, error) {
	return atomfs.db.CreateMolecule(name, atoms)
}

// CopyMolecule simply duplicates a molecule's configuration under a new name.
// This is equivalent to a "snapshot" operation under other filesystems.
func (atomfs *Instance) CopyMolecule(dest string, src string) (types.Molecule, error) {
	mol, err := atomfs.db.GetMolecule(src)
	if err != nil {
		return types.Molecule{}, err
	}

	return atomfs.db.CreateMolecule(dest, mol.Atoms)
}

func (atomfs *Instance) DeleteMolecule(name string) error {
	mol, err := atomfs.db.GetMolecule(name)
	if err != nil {
		return err
	}

	return atomfs.db.DeleteThing(mol.ID, "molecule")
}

func (atomfs *Instance) RenameMolecule(old, new_ string) error {
	mol, err := atomfs.db.GetMolecule(old)
	if err != nil {
		return err
	}
	return atomfs.db.RenameThing(mol.ID, "molecule", new_)
}

func (atomfs *Instance) CreateMoleculeFromOCITag(oci casext.Engine, name string) (types.Molecule, error) {
	man, err := stacker.LookupManifest(oci, name)
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

		atom, err := atomfs.CreateAtomFromOCIBlob(layer)
		if err != nil {
			return types.Molecule{}, err
		}

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

	return atomfs.db.CreateMolecule(name, atoms)
}

func (atomfs *Instance) GetMolecule(name string) (types.Molecule, error) {
	return atomfs.db.GetMolecule(name)
}
