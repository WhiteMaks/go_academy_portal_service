package util

import "github.com/spf13/viper"

type Config struct {
	Microservice Microservice `mapstructure:"microservice"`
	Database     Database     `mapstructure:"database"`
}

type Microservice struct {
	Port int `mapstructure:"port"`
}

type Database struct {
	ConnectionString string `mapstructure:"connection_string"`
}

func LoadConfig(path string) (Config, error) {
	var result = Config{}

	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		return result, err
	}

	err = viper.Unmarshal(&result)
	if err != nil {
		return result, err
	}

	return result, nil
}
