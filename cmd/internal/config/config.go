package config

import (
	"rest-api/pkg/logging"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		IsDebug *bool   `yaml:"is_debug"`
		Listen  Listen  `yaml:"listen"`
		MongoDB MongoDB `yaml:"mongodb"`
	}

	Listen struct {
		Type   string `yaml:"type"`
		BindIp string `yaml:"bind_ip"`
		Port   string `yaml:"port"`
	}

	MongoDB struct {
		Host       string `yaml:"host"`
		Port       string `yaml:"port"`
		Database   string `yaml:"database"`
		AuthDB     string `yaml:"auth_db"`
		Username   string `yaml:"username"`
		Password   string `yaml:"password"`
		Collection string `yaml:"collection"`
	}
)

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("Read appplication configuration")
		instance = &Config{}
		if err := cleanenv.ReadConfig("../../config.yaml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Info(help)
			logger.Fatal(err)
		}
	})

	return instance
}
