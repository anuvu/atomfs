package atomfs

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/anuvu/atomfs/types"
)

func TestRename(t *testing.T) {
	dir, err := ioutil.TempDir("", "atomfs-rename-")
	if err != nil {
		t.Fatalf("couldn't make tempdir %s", err)
	}
	defer os.RemoveAll(dir)

	atomfs, err := New(types.Config{Path: dir})
	if err != nil {
		t.Fatalf("couldn't open atomfs %s", err)
	}

	mol1, err := atomfs.CreateMolecule("foo", nil)
	if err != nil {
		t.Fatalf("couldn't create molecule %s", err)
	}

	err = atomfs.RenameMolecule("foo", "bar")
	if err != nil {
		t.Fatalf("couldn't rename molecule %s", err)
	}

	mol2, err := atomfs.GetMolecule("bar")
	if err != nil {
		t.Fatalf("couldn't get molecule %s", err)
	}

	if mol1.ID != mol2.ID {
		t.Fatalf("molecule ids changed after rename")
	}
}
