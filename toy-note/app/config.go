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
	viper.SetConfigType("env") // 可以使用 JSON，XML等其它格式

	// 自动覆盖环境变量值
	viper.AutomaticEnv()

	// 开始读取配置值
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// 转换为 Config 结构体的变量
	err = viper.Unmarshal(&config)
	return
}
