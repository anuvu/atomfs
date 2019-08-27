package atomfs

import (
	"github.com/anuvu/atomfs/mount"
)

func (atomfs *Instance) Mount(molecule string, target string, writable bool) error {
	mol, err := atomfs.db.GetMolecule(molecule)
	if err != nil {
		return err
	}

	ovl, err := mount.NewOverlay(atomfs.config, mol)
	if err != nil {
		return err
	}

	return ovl.Mount(target, writable)
}

func (atomfs *Instance) Umount(target string) error {
	return mount.Umount(atomfs.config, target)
}
