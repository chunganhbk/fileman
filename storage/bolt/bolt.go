package bolt

import (
	"github.com/asdine/storm"
	"github.com/raedahgroup/fileman/config"
	"github.com/raedahgroup/fileman/storage"
	"github.com/raedahgroup/fileman/users"
)

// NewStorage creates a storage.Storage based on Bolt DB.
func NewStorage(db *storm.DB) (*storage.Storage, error) {
	users := users.NewStorage(usersBackend{db: db})
	config := config.NewStorage(configBackend{db: db})
	err := save(db, "version", 1)
	if err != nil {
		return nil, err
	}

	return &storage.Storage{
		Users:    users,
		Config: config,
	}, nil
}
