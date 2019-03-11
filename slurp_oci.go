package atomfs

import (
	"context"
	"fmt"
	"io"

	"github.com/anuvu/atomfs/types"
	"github.com/anuvu/stacker"
	"github.com/openSUSE/umoci"
	"github.com/openSUSE/umoci/oci/casext"
	ispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
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
	for i, l := range man.Layers {
		layer, err := oci.FromDescriptor(context.Background(), l)
		if err != nil {
			return err
		}
		defer layer.Close()

		atomType := types.TarAtom
		switch layer.Descriptor.MediaType {
		case ispec.MediaTypeImageLayer:
			fallthrough
		case ispec.MediaTypeImageLayerGzip:
			fallthrough
		case ispec.MediaTypeImageLayerNonDistributable:
			fallthrough
		case ispec.MediaTypeImageLayerNonDistributableGzip:
			atomType = types.TarAtom
		default:
			return errors.Errorf("unknown media type: %s", layer.Descriptor.MediaType)
		}

		name := fmt.Sprintf("%s-%d", name, i)
		atom, err := atomfs.db.CreateAtom(name, atomType, layer.Data.(io.Reader))
		if err != nil {
			return err
		}

		atoms = append(atoms, atom)
	}

	_, err = atomfs.db.CreateMolecule(name, atoms)
	return err
}
