package config

import (
	"strings"

	"github.com/spf13/viper"
)

var config *viper.Viper

// Init is used to initialize the configuration for this intance of the program.
// The current directory "config" is searched for a json file whose name matches
// {env}.json, where {env} is the environment the program is running in.
func Init(env string) {
	config = viper.New()
	config.SetConfigType("json")
	config.SetConfigName(env)
	config.AddConfigPath("config/")

	config.SetEnvPrefix("followrs")
	config.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	config.SetEnvKeyReplacer(replacer)

	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}
}

// GetConfig returns the config for this instance of the program.
func GetConfig() *viper.Viper {
	return config
}
