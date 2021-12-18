package app

import (
	"github.com/spf13/viper"
)

type Config struct {
	PG_HOST    string
	PG_PORT    int
	PG_USER    string
	PG_PASS    string
	PG_DB      string
	MONGO_HOST string
	MONGO_PORT int
	MONGO_USER string
	MONGO_PASS string
	MONGO_DB   string
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// auto-override environment config
	viper.AutomaticEnv()

	// read config
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// unmarshal to Config struct
	err = viper.Unmarshal(&config)
	return
}
