package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
)

func TestConfig_EnvironmentOverride(t *testing.T) {

	assert := assert.New(t)

	err := os.Setenv("CONFIG_PATH", "../../configs/config.env")
	assert.NoError(err)

	override := Config{
		DbConnectionString: "newConnectionString",
		Port:               1234,
		Mode:               "release",
	}

	_ = os.Setenv("DB_CONNECTION_STRING", override.DbConnectionString)
	_ = os.Setenv("PORT", strconv.Itoa(override.Port))
	_ = os.Setenv("MODE", override.Mode)

	cfg := Get()

	assert.Equal(override.DbConnectionString, cfg.DbConnectionString)
	assert.Equal(override.Port, cfg.Port)
	assert.Equal(override.Mode, cfg.Mode)
}
