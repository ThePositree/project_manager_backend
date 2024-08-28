package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type config struct {
	MongoURI          string `json:"mongo_uri"`
	Database          string `json:"database"`
	UserCollection    string `json:"user_collection"`
	BillingCollection string `json:"billing_collection"`
	AdminPassword     string `json:"admin_password"`
	HttpPort          int    `json:"http_port"`
}

func New() (config, error) {
	bytes, err := os.ReadFile("config.json")
	if err != nil {
		return config{}, fmt.Errorf("read file: %w", err)
	}
	var result config
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return config{}, fmt.Errorf("unmarshal json: %w", err)
	}
	return result, nil
}
