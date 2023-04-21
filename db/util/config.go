package util

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

// Config store all configuration of the application
// The values are read by viper from a config file or environment variable
type Config struct {
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	RedisAddress         string         `mapstructure:"REDIS_ADDRESS"`
	TokenDuration        time.Duration `mapstructure:"TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

// LoadConfig reads configurations from file or environment
func LoadConfig(path string) (c Config, err error) {

	viper.AddConfigPath(path)

	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil { // Handle errors reading the config file
		fmt.Errorf("fatal error config file: %w", err)
		return
	}

	err = viper.Unmarshal(&c)
	return

}
