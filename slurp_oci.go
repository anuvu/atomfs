package atomfs

import (
	"context"

	"github.com/opencontainers/umoci"
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
		if _, err := atomfs.CreateMoleculeFromOCITag(oci, t); err != nil {
			return err
		}
	}

	return nil
}
