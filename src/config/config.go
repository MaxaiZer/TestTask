package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"strings"
)

type Config struct {
	Port int       `mapstructure:"port"`
	Mode string    `mapstructure:"mode"`
	JWT  JWTConfig `mapstructure:"jwt"`
	DB   DBConfig  `mapstructure:"db"`
}

type JWTConfig struct {
	AccessLifetime  int    `mapstructure:"access_lifetime"`
	RefreshLifetime int    `mapstructure:"refresh_lifetime"`
	SecretKey       string `mapstructure:"secret_key"`
	Issuer          string `mapstructure:"issuer"`
	Audience        string `mapstructure:"audience"`
}

type DBConfig struct {
	ConnectionString string `mapstructure:"connection_string"`
}

var configFile = "config.yaml"

func Get() *Config {

	config, err := loadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func loadConfig(file string) (*Config, error) {

	viper.SetConfigFile(file)
	viper.AutomaticEnv()

	_ = viper.BindEnv("jwt.secret_key", "JWT_SECRET_KEY")
	_ = viper.BindEnv("jwt.access_lifetime", "JWT_ACCESS_LIFETIME")
	_ = viper.BindEnv("jwt.refresh_lifetime", "JWT_REFRESH_LIFETIME")
	_ = viper.BindEnv("jwt.issuer", "JWT_ISSUER")
	_ = viper.BindEnv("jwt.audience", "JWT_AUDIENCE")
	_ = viper.BindEnv("db.connection_string", "DB_CONNECTION_STRING")

	viper.SetDefault("PORT", 8080)
	viper.SetDefault("MODE", "release")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	err := config.DB.Validate()
	if err != nil {
		return nil, err
	}

	err = config.JWT.Validate()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (config DBConfig) Validate() error {
	if config.ConnectionString == "" {
		return fmt.Errorf("missing variable: db connection string")
	}
	return nil
}

func (config JWTConfig) Validate() error {

	var missingFields []string

	if config.SecretKey == "" {
		missingFields = append(missingFields, "secret key")
	}

	if config.Issuer == "" {
		missingFields = append(missingFields, "issuer")
	}

	if config.Audience == "" {
		missingFields = append(missingFields, "audience")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required variables: %s", strings.Join(missingFields, ", "))
	}

	if config.AccessLifetime <= 0 {
		return fmt.Errorf("invalid access lifetime: must be greater than 0")
	}

	if config.RefreshLifetime <= 0 {
		return fmt.Errorf("invalid refresh lifetime: must be greater than 0")
	}

	return nil
}
