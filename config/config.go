
package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var State ConfigState

type ConfigState struct {
	Debug        bool   `yaml:"debug"`
	Port         string `yaml:"port"`
	APIVersion   string `yaml:"api_version"`
	JWTKEY       string `yaml:"jwt_key"`
	DatabasePath string `yaml:"database_path"`
	BaseURL string `yaml:"base_url"`
	RootPath string `yaml:"root_path"`

}
func Load(filePath string) error {
	configFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(configFile, &State)
	return err
}

