package storage

import (
	"github.com/raedahgroup/fileman/config"
	"github.com/raedahgroup/fileman/users"
)

// Storage is a storage powered by a Backend whih makes the neccessary
// verifications when fetching and saving data to ensure consistency.
type Storage struct {
	Users    *users.Storage
	Config   *config.Storage
}
