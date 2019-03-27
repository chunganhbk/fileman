package bolt

import (
	"github.com/asdine/storm"
	"github.com/raedahgroup/fileman/config"
)

type configBackend struct {
	db *storm.DB
}
func (s configBackend) SaveServer(config *config.Server) error {
	return save(s.db, "server", config)
}

func (s configBackend) GetServer() (*config.Server, error) {
	server := &config.Server{}
	return server, get(s.db, "server", server)
}
