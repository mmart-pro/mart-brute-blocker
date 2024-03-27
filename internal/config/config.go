package config

import (
	"encoding/json"
	"net"
	"os"
)

type LoggerConfig struct {
	Level   string `json:"level"`
	LogFile string `json:"logfile"`
}

type GrpcConfig struct {
	GrpcHost string `json:"host"`
	GrpcPort string `json:"port"`
}

func (c GrpcConfig) GetEndpoint() string {
	return net.JoinHostPort(c.GrpcHost, c.GrpcPort)
}

type ServiceConfig struct {
	MaxPerMinForLogin    int `json:"maxPerMinForLogin"`
	MaxPerMinForPassword int `json:"maxPerMinForPassword"`
	MaxPerMinForIP       int `json:"maxPerMinForIp"`
}

type StorageConfig struct {
	UseDB    bool   `json:"useDb"`
	Host     string `json:"host"`
	Port     int16  `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type APIConfig struct {
	LoggerConfig  LoggerConfig  `json:"logger"`
	GrpcConfig    GrpcConfig    `json:"grpc"`
	ServiceConfig ServiceConfig `json:"service"`
	StorageConfig StorageConfig `json:"storage"`
}

func NewAPIConfig(filename string) (APIConfig, error) {
	var cfg APIConfig
	return cfg, newConfig(filename, &cfg)
}

type CliConfig struct {
	GrpcConfig GrpcConfig `json:"grpc"`
}

func NewCliConfig(filename string) (CliConfig, error) {
	var cfg CliConfig
	return cfg, newConfig(filename, &cfg)
}

func newConfig(filename string, dest interface{}) error {
	arr, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(arr, dest)
	if err != nil {
		return err
	}

	return nil
}
