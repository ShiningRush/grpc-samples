package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"gitlab.followme.com/FollowmeGo/golib/grpc/utils"
)

type configOption struct {
	dir  string
	name string
}

type OptAction func(*configOption)

func InitConfig(opts ...OptAction) {
	opt := initOpts(opts)

	initBaseConfig(opt)
	mergeEnvConfig(utils.GetEnvOrDefault("RUN_ENV", "development"))
}

func initOpts(opts []OptAction) *configOption {
	opt := &configOption{
		dir:  "./configs",
		name: "config",
	}

	for _, v := range opts {
		v(opt)
	}

	return opt
}

func initBaseConfig(opt *configOption) {
	viper.AddConfigPath(opt.dir) // optionally look for config in the working directory
	viper.SetConfigName(opt.name)
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

func GetViper() *viper.Viper {
	return viper.GetViper()
}
