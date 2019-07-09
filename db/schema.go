package db

import (
	"database/sql"
	"fmt"

	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

var Schema string = `
CREATE TABLE IF NOT EXISTS schema (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	version INTEGER NOT NULL,
	updated DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS atoms (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	name TEXT NOT NULL,
	hash TEXT NOT NULL,
	type TEXT NOT NULL,
	UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS molecules (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	name TEXT NOT NULL,
	UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS molecule_atoms (
	id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	molecule_id INTEGER NOT NULL,
	atom_id INTEGER NOT NULL,
	FOREIGN KEY (molecule_id) REFERENCES molecules (id) ON DELETE CASCADE,
	-- Note: we explicitly do not want ON DELETE CASCADE here. If we
	-- automatically delete unused atoms, we won't know to delete them from
	-- the FS.
	FOREIGN KEY (atom_id) REFERENCES atoms (id)
);
`

func init() {
	sql.Register("sqlite3_with_fk", &sqlite3.SQLiteDriver{ConnectHook: sqliteEnableForeignKeys})
}

func sqliteEnableForeignKeys(conn *sqlite3.SQLiteConn) error {
	_, err := conn.Exec("PRAGMA foreign_keys=ON;", nil)
	return err
}

func openSqlite(path string) (*sql.DB, error) {
	openPath := fmt.Sprintf("%s?_busy_timeout=5&_txlock=exclusive", path)
	db, err := sql.Open("sqlite3_with_fk", openPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(Schema)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't create schema")
	}

	return db, nil
}
