package config

import (
	"errors"
	"fmt"

	"github.com/boldlogic/packages/commonconfig"
	"github.com/boldlogic/packages/dbconfig"
	logger "github.com/boldlogic/packages/logger/zaplog"
	"github.com/boldlogic/packages/transport/httpserver"
)

type Config struct {
	HTTP     httpserver.ServerConfig `yaml:"http" json:"http"`
	Logger   logger.Config           `yaml:"logger" json:"logger"`
	Database dbconfig.DBConfig       `yaml:"database" json:"database"`
}

func Load(configPath string) (*Config, error) {

	cfg, err := commonconfig.DecodeConfigStrict[Config](configPath)
	if err != nil {
		return nil, err
	}

	cfg.applyDefaults()
	errs := cfg.validate()
	if err := errors.Join(errs...); err != nil {
		return nil, fmt.Errorf("некорректный конфиг: %w", err)
	}
	return &cfg, nil
}

func (c *Config) validate() []error {
	var errs []error

	logErr := c.Logger.Validate()
	if len(logErr) > 0 {
		errs = append(errs, logErr...)
	}

	dbErrs := c.Database.Validate()
	if len(dbErrs) > 0 {
		errs = append(errs, dbErrs...)
	}

	srvErrs := c.HTTP.Validate()
	if len(srvErrs) > 0 {
		errs = append(errs, srvErrs...)
	}
	return errs
}

func (c *Config) applyDefaults() {
	c.Database.ApplyDefaults()
	c.HTTP.ApplyDefaults()

}
