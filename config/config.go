package config

import (
	"encoding/json"
	"os"

	"github.com/thedanisaur/jfl_platform/types"
)

func loadConfig(config_path string) (types.Config, error) {
	var config types.Config
	config_file, err := os.Open(config_path)
	if err != nil {
		return config, err
	}
	defer config_file.Close()
	jsonParser := json.NewDecoder(config_file)
	jsonParser.Decode(&config)
	return config, nil
}
