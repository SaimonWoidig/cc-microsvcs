package config

import "github.com/spf13/viper"

func NewViper() *viper.Viper {
	v := viper.NewWithOptions()

	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.SetConfigName("config.yaml")
	v.SetConfigType("yaml")

	return v
}
