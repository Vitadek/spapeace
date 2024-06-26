package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type Config struct {
	AnsibleVaultLocation string `json:"ansible_vault_location"`
	AnsibleVaultPassword string `json:"ansible_vault_password"`
	HostLimitations      string `json:"host_limitations"`
}

var configFilePath = "config.json"
var currentConfig Config

func LoadConfig() error {
	bytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return errors.New("failed to read config file: " + err.Error())
	}

	err = json.Unmarshal(bytes, &currentConfig)
	if err != nil {
		return errors.New("failed to parse config: " + err.Error())
	}

	return nil
}

func SaveConfig() error {
	bytes, err := json.MarshalIndent(currentConfig, "", "  ")
	if err != nil {
		return errors.New("failed to marshal config: " + err.Error())
	}

	err = ioutil.WriteFile(configFilePath, bytes, 0644)
	if err != nil {
		return errors.New("failed to write config file: " + err.Error())
	}

	return nil
}

func GetConfig() *Config {
	return &currentConfig
}
