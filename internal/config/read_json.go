package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadJsonFile() Config {
	cfgFilePath, err := getConfigFilePath()
	if err != nil {
		fmt.Println("Config file not found in home directory, please create a config file!")
	}

	var cfg Config
	jsonFile, err := os.ReadFile(cfgFilePath)
	if err != nil {
		fmt.Println(err)
		return Config{}
	}
	if err := json.Unmarshal(jsonFile, &cfg); err != nil {
		fmt.Println(err)
		return Config{}
	}

	return cfg
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return homeDir + "/" + configFileName, nil
}
