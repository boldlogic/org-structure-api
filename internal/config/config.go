package config

import (
	"github.com/boldlogic/packages/commonconfig"
	logger "github.com/boldlogic/packages/logger/zaplog"
	"github.com/boldlogic/packages/transport/httpserver"
)

type Config struct {
	HTTP     httpserver.ServerConfig `yaml:"http" json:"http"`
	Logger   logger.Config           `yaml:"logger" json:"logger"`
	Database DBConfig                `yaml:"database" json:"database"`
}

func Load(configPath string) (*Config, error) {

	cfg, err := commonconfig.DecodeConfigStrict[Config](configPath)
	if err != nil {
		return nil, err
	}

	cfg.HTTP.ApplyDefaults()

	return &cfg, nil
}
