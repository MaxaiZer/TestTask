package config

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

const (
	DebugMode   = "debug"
	ReleaseMode = "release"
)

type Config struct {
	Port               int    `mapstructure:"PORT"`
	Mode               string `mapstructure:"MODE"`
	DbConnectionString string `mapstructure:"DB_CONNECTION_STRING"`
}

var configFile = "./configs/config.env"

func Get() *Config {

	config, err := loadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func loadConfig(file string) (*Config, error) {

	if os.Getenv("CONFIG_PATH") != "" {
		file = os.Getenv("CONFIG_PATH")
	}

	viper.SetConfigFile(file)
	viper.AutomaticEnv()

	viper.SetDefault("PORT", 8080)
	viper.SetDefault("MODE", ReleaseMode)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	err := config.validate()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c Config) validate() error {

	var errs []error

	if c.DbConnectionString == "" {
		errs = append(errs, fmt.Errorf("missing variable DbConnectionString"))
	}

	if c.Port <= 0 {
		errs = append(errs, fmt.Errorf("invalid port: %d", c.Port))
	}

	if c.Mode != ReleaseMode && c.Mode != DebugMode {
		errs = append(errs, fmt.Errorf("invalid mode: %s", c.Mode))
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors occurred: %w", errors.Join(errs...))
	}

	return nil
}
