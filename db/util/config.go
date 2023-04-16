package util

import (
	"github.com/spf13/viper"
)

// Config store all configuration of the application
// The values are read by viper from a config file or environment variable
type Config struct {
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	DBSource      string `mapstructure:"DB_SOURCE"`
}

// LoadConfig reads configurations from file or environment
func LoadConfig(path string) (c Config, err error) {

	viper.AddConfigPath(path)

	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil { // Handle errors reading the config file

		return
	}

	err = viper.Unmarshal(&c)
	return

}
