package atomfs

import (
	"github.com/anuvu/atomfs/types"
)

func (atomfs *Instance) Mount(molecule string, target string, writable bool) error {
	mol, err := atomfs.db.GetMolecule(molecule)
	if err != nil {
		return err
	}

	return mol.Mount(target, writable)
}

func (atomfs *Instance) Umount(target string) error {
	return types.Umount(atomfs.config, target)
}
