package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	Teams       []Team                `json:"teams"`
	Permissions map[string]Permission `json:"permissions"`
	ApiKey      string                `json:"apikey"`
}

func NewConfigFromBytes(bytes []byte) (c *Config, err error) {
	c = &Config{}
	err = json.Unmarshal(bytes, c)
	return
}
func NewConfigFromFile(filename string) (c *Config, err error) {
	jsonFile, err := os.Open(filename)
	defer func() {
		_ = jsonFile.Close()
	}()
	if err != nil {
		return
	}
	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return
	}
	c, err = NewConfigFromBytes(bytes)
	return
}
