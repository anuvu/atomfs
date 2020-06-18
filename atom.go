package atomfs

import (
	"io"

	"github.com/anuvu/atomfs/types"
	"github.com/opencontainers/umoci/oci/casext"
	ispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

func (atomfs *Instance) GetAtoms() ([]types.Atom, error) {
	return atomfs.db.GetAtoms()
}

func (atomfs *Instance) CreateAtom(name string, atomType types.AtomType, content io.Reader) (types.Atom, error) {
	return atomfs.db.CreateAtom(name, atomType, content)
}

func (atomfs *Instance) CreateAtomFromOCIBlob(blob *casext.Blob) (types.Atom, error) {
	atomType := types.TarAtom
	switch blob.Descriptor.MediaType {
	case ispec.MediaTypeImageLayer:
		fallthrough
	case ispec.MediaTypeImageLayerGzip:
		fallthrough
	case ispec.MediaTypeImageLayerNonDistributable:
		fallthrough
	case ispec.MediaTypeImageLayerNonDistributableGzip:
		atomType = types.TarAtom
	// stolen from stacker:base.go
	case "application/vnd.oci.image.layer.squashfs":
		atomType = types.SquashfsAtom
	default:
		return types.Atom{}, errors.Errorf("unknown media type: %s", blob.Descriptor.MediaType)
	}

	return atomfs.db.CreateAtom(blob.Descriptor.Digest.Encoded(), atomType, blob.Data.(io.Reader))
}

func (atomfs *Instance) GetAtomsByHash() (map[string]types.Atom, error) {
	atomsList, err := atomfs.GetAtoms()
	if err != nil {
		return nil, err
	}

	atoms := map[string]types.Atom{}
	for _, atom := range atomsList {
		atoms[atom.Hash] = atom
	}

	return atoms, nil
}
