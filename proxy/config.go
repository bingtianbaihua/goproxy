package proxy

import (
	"bufio"
	"encoding/json"
	"os"
)

type Config struct {
	Port string `json:"port"`
	Auth bool   `json:"auth"`

	User map[string]string `json:"user"`
}

func (c *Config) GetConfig(filename string) error {

	configFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer configFile.Close()

	br := bufio.NewReader(configFile)
	err = json.NewDecoder(br).Decode(c)
	if err != nil {
		return err
	}
	return nil
}
