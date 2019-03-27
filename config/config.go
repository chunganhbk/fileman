
package config

import (
	"math/rand"
	"strings"
)

type Server struct {
	JWTKEY       string `yaml:"jwt_key"`
	RootPath string `yaml:"root_path"`
	BaseURL string `yaml:"root_path"`

}

// Clean cleans any variables that might need cleaning.
func (s *Server) Clean() {
	s.BaseURL = strings.TrimSuffix(s.BaseURL, "/")
}

// GenerateKey generates a key of 256 bits.
func GenerateKey() ([]byte, error) {
	b := make([]byte, 64)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}
