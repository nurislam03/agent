package config

import (
	"github.com/sirupsen/logrus"
	"sync"

	// viper/remote
	_ "github.com/spf13/viper/remote"
)

// Config ...
type Config struct {
}

func loadConfig() {
	logrus.Info("Loading configurations...")
	logrus.Info("Configurations loaded")
}

var config *Config
var configOnce = &sync.Once{}

// NewConfig ...
func NewConfig() *Config {
	configOnce.Do(func() {
		//loadConfig()
		config = &Config{
		}
	})
	return config
}
