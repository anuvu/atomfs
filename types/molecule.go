package types

type Molecule struct {
	ID   int64
	Name string
	// Atoms is the list of atoms in this Molecule. The first element in
	// this list is the top most layer in the overlayfs.
	Atoms []Atom

	config Config
}

func NewMolecule(config Config) Molecule {
	return Molecule{config: config}
}


