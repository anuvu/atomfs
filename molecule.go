package atomfs

import (
	"github.com/anuvu/atomfs/types"
)

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
