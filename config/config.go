package config

import (
	"log"

	"github.com/spf13/viper"
)

var config *viper.Viper

func Init(env string) {
	config = viper.New()
	config.SetConfigType("json")
	config.SetConfigName(env)
	config.AddConfigPath("../config/")
	config.AddConfigPath("config/")
	var err error = config.ReadInConfig()
	if err != nil {
		log.Print(err)
		log.Fatal("error parsing the configuration file")

	}
}

func GetConfig() *viper.Viper {
	return config
}
