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

	requiredFields := []string{
		"jwt.secret_key",
		"db.connection_string",
		"jwt.access_lifetime",
		"jwt.refresh_lifetime",
		"jwt.issuer",
		"jwt.audience",
	}

	var missingFields []string

	for _, field := range requiredFields {
		if !viper.IsSet(field) || viper.GetString(field) == "" {
			missingFields = append(missingFields, field)
		}
	}

	if len(missingFields) > 0 {
		var fieldsStr = strings.Join(missingFields, ", ")
		return nil, fmt.Errorf("missing or empty required fields in config: %s", fieldsStr)
	}

	return &config, nil
}
