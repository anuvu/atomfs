package db

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/anuvu/atomfs/types"
	"github.com/pkg/errors"
)

type AtomfsDB struct {
	// Expose the DB; although nobody should use it because the helper
	// methods should be ok, you never know...
	DB     *sql.DB
	config types.Config
}

func New(config types.DBBasedConfig) (*AtomfsDB, error) {
	p := config.RelativePath("atomfs.db")
	db, err := openSqlite(p)
	if err != nil {
		return nil, errors.Wrapf(err, "error opening db at %s", p)
	}

	return &AtomfsDB{db, config}, nil
}

func (db *AtomfsDB) Close() error {
	return db.DB.Close()
}

func (db *AtomfsDB) CreateAtom(name string, atomType types.AtomType, content io.Reader) (types.Atom, error) {
	f, err := ioutil.TempFile(db.config.AtomsPath(), "create-atom-")
	if err != nil {
		return types.Atom{}, errors.Wrapf(err, "error creating tempfile for atom import")
	}
	defer f.Close()

	h := sha256.New()
	w := io.MultiWriter(h, f)

	_, err = io.Copy(w, content)
	if err != nil {
		return types.Atom{}, errors.Wrapf(err, "error copying data to atom import file")
	}

	hash := fmt.Sprintf("%x", h.Sum(nil))
	f.Close()
	err = os.Rename(f.Name(), db.config.AtomsPath(hash))
	if err != nil {
		return types.Atom{}, errors.Wrapf(err, "error renaming atom import file")
	}

	stmt, err := db.DB.Prepare("INSERT INTO atoms (name, hash, type) VALUES (?, ?, ?)")
	if err != nil {
		return types.Atom{}, errors.Wrapf(err, "error preparing atom insert")
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, hash, atomType)
	if err != nil {
		return types.Atom{}, errors.Wrapf(err, "error inserting atom")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return types.Atom{}, errors.Wrapf(err, "error getting last atom insert id")
	}

	return types.Atom{id, name, hash, atomType}, nil
}

func (db *AtomfsDB) getAtoms(rows *sql.Rows) ([]types.Atom, error) {
	atoms := []types.Atom{}
	for rows.Next() {
		atom := types.Atom{}
		err := rows.Scan(&atom.ID, &atom.Name, &atom.Hash, &atom.Type)
		if err != nil {
			return nil, errors.Wrapf(err, "error scanning row")
		}
		atoms = append(atoms, atom)
	}

	return atoms, nil
}

func (db *AtomfsDB) GetAtoms() ([]types.Atom, error) {
	rows, err := db.DB.Query("SELECT id, name, hash, type FROM atoms")
	if err != nil {
		return nil, errors.Wrapf(err, "error querying atoms")
	}
	defer rows.Close()

	return db.getAtoms(rows)
}

func (db *AtomfsDB) CreateMolecule(name string, atoms []types.Atom) (types.Molecule, error) {
	stmt, err := db.DB.Prepare("INSERT INTO molecules (name) VALUES (?)")
	if err != nil {
		return types.Molecule{}, errors.Wrapf(err, "error preparing molecule insert")
	}

	result, err := stmt.Exec(name)
	stmt.Close()
	if err != nil {
		return types.Molecule{}, errors.Wrapf(err, "error inserting molecule")
	}

	stmt, err = db.DB.Prepare("INSERT INTO molecule_atoms (molecule_id, atom_id) VALUES (?, ?)")
	if err != nil {
		return types.Molecule{}, err
	}
	defer stmt.Close()

	id, err := result.LastInsertId()
	if err != nil {
		return types.Molecule{}, err
	}

	for _, a := range atoms {
		_, err = stmt.Exec(id, a.ID)
		if err != nil {
			// TODO: cleanup?
			return types.Molecule{}, err
		}
	}

	mol := types.NewMolecule(db.config)
	mol.ID = id
	mol.Name = name
	mol.Atoms = atoms
	return mol, nil
}

func (db *AtomfsDB) GetMolecule(name string) (types.Molecule, error) {
	rows, err := db.DB.Query("SELECT id, name FROM molecules WHERE name=?", name)
	if err != nil {
		return types.Molecule{}, err
	}
	defer rows.Close()

	mol := types.NewMolecule(db.config)
	found := false
	for rows.Next() {
		err = rows.Scan(&mol.ID, &mol.Name)
		if err != nil {
			return types.Molecule{}, err
		}
		found = true
	}

	if !found {
		return types.Molecule{}, errors.Errorf("molecule %s not found", name)
	}

	rows, err = db.DB.Query(`
		SELECT atoms.id, atoms.name, atoms.hash, atoms.type
		FROM atoms JOIN molecule_atoms ON atoms.id = molecule_atoms.atom_id
		WHERE molecule_atoms.molecule_id = ?
		ORDER BY molecule_atoms.id ASC`, mol.ID)
	if err != nil {
		return types.Molecule{}, err
	}
	defer rows.Close()

	mol.Atoms, err = db.getAtoms(rows)
	if err != nil {
		return types.Molecule{}, err
	}

	return mol, nil
}

func (db *AtomfsDB) GetUnusedAtoms() ([]types.Atom, error) {
	rows, err := db.DB.Query(`
		SELECT atoms.id, atoms.name, atoms.hash, atoms.type
		FROM atoms
		WHERE atoms.id not in (
			SELECT atoms.id
			FROM atoms JOIN molecule_atoms ON atoms.id = molecule_atoms.atom_id
		)`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return db.getAtoms(rows)
}

func (db *AtomfsDB) DeleteThing(id int64, table string) error {
	_, err := db.DB.Exec(fmt.Sprintf("DELETE FROM %ss WHERE id = ?", table), id)
	return err
}

func (db *AtomfsDB) RenameThing(id int64, table string, newName string) error {
	_, err := db.DB.Exec(fmt.Sprintf("UPDATE %ss SET name = ? WHERE id = ?", table), newName, id)
	return err
}
