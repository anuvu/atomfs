package atomfs

import (
	"os"

	"github.com/anuvu/atomfs/db"
	"github.com/anuvu/atomfs/types"
)

type Instance struct {
	config types.Config
	db     *db.AtomfsDB
}

func New(config types.Config) (*Instance, error) {
	if err := os.MkdirAll(config.Path, 0755); err != nil {
		if !os.IsExist(err) {
			return nil, err
		}
	}

	db, err := db.New(config)
	if err != nil {
		return nil, err
	}

	return &Instance{config: config, db: db}, nil
}

func (atomfs *Instance) Close() error {
	return atomfs.db.Close()
}
