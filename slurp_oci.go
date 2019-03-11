package atomfs

import (
	"context"

	"github.com/anuvu/atomfs/types"
	"github.com/anuvu/stacker"
	"github.com/openSUSE/umoci"
	"github.com/openSUSE/umoci/oci/casext"
)

func (atomfs *Instance) SlurpOCI(location string) error {
	oci, err := umoci.OpenLayout(location)
	if err != nil {
		return err
	}
	defer oci.Close()

	tags, err := oci.ListReferences(context.Background())
	if err != nil {
		return err
	}

	for _, t := range tags {
		if err := atomfs.slurpTag(oci, t); err != nil {
			return err
		}
	}

	return nil
}

func (atomfs *Instance) slurpTag(oci casext.Engine, name string) error {
	man, err := stacker.LookupManifest(oci, name)
	if err != nil {
		return err
	}

	atoms := []types.Atom{}
	for _, l := range man.Layers {
		layer, err := oci.FromDescriptor(context.Background(), l)
		if err != nil {
			return err
		}
		defer layer.Close()

		atom, err := atomfs.CreateAtomFromOCIBlob(layer)
		if err != nil {
			return err
		}

		atoms = append(atoms, atom)
	}

	_, err = atomfs.db.CreateMolecule(name, atoms)
	return err
}
