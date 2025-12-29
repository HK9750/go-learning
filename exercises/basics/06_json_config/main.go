package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	AppName    string   `json:"app_name"`
	Port       int      `json:"port"`
	Debug      bool     `json:"debug"`
	AllowedIps []string `json:"allowed_ips"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config json: %w", err)
	}
	defer file.Close()

	var config Config

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config json: %w", err)
	}

	if config.AppName == "" {
		return nil, fmt.Errorf("validation error: app_name cannot be empty")
	}

	if config.Port < 1 || config.Port > 65535 {
		return nil, fmt.Errorf("validation error: port must be between 1 and 65535")
	}

	return &config, nil
}

func main() {
	config, err := LoadConfig("./exercises/06_json_config/config.json")
	if err != nil {
		fmt.Println("Error in main", err)
	}
	fmt.Printf("%T", config)
}
