package config

import (
	"fmt"
	"strings"

	"github.com/shiningrush/grpc-samples/pkg/grpc/utils"
	"github.com/spf13/viper"
)

func init() {
	initBaseConfig()
	mergeEnvConfig(utils.GetEnvOrDefault("RUN_ENV", "development"))
}

func initBaseConfig() {
	viper.AddConfigPath("./configs") // optionally look for config in the working directory
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

}

func mergeEnvConfig(env string) {
	viper.SetConfigName("config." + strings.ToLower(env))
	if err := viper.MergeInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

func GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}

func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

func Get(key string) interface{} {
	return viper.Get(key)
}

func Set(key string, value interface{}) {
	viper.Set(key, value)
}

func IsSet(key string) bool {
	return viper.IsSet(key)
}
