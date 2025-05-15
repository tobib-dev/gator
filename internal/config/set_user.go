package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func SetUser(userName string) {
	cfgFilePath, err := getConfigFilePath()
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonFile, err := os.ReadFile(cfgFilePath)
	if err != nil {
		fmt.Println(err)
	}

	var gatorCfg Config
	if err := json.Unmarshal(jsonFile, &gatorCfg); err != nil {
		fmt.Println(err)
	}

	gatorCfg.CurrentUserName = userName
	err = write(gatorCfg)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

func write(cfg Config) error {
	cfgFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	os.WriteFile(cfgFilePath, jsonData, 0644)
	return nil
}
