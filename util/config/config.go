package config

import (
	"flag"
	"fmt"
	"github.com/spf13/cast"
	"github.com/spf13/viper" // 自定义包名，避免与内置 viper 实例冲突
)

var v *viper.Viper

// 初始化配置信息，完成对环境变量以及 conf 信息的加载
func init() {
	// 使用独立的实例。防止外部直接调用 viper 标准实例
	v = viper.New()
	v.SetConfigType("toml") // "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"
	v.AddConfigPath(".")
	configFlag := flag.String("config", "config.toml", "path to config file")
	// 将命令行标志的值设置为 Viper 配置实例的属性
	v.SetConfigFile(*configFlag)
	if err := v.ReadInConfig(); err != nil {
		fmt.Println("ReadInConfig err", err)
	}
	v.WatchConfig()
}

// Env 读取环境变量，支持默认值
func Env(envName string, defaultValue ...any) any {
	if len(defaultValue) > 0 {
		return internalGet(envName, defaultValue[0])
	}
	return internalGet(envName)
}

// Get 获取配置项
// 第一个参数 path 允许使用点式获取，如：app.name
// 第二个参数允许传参默认值
func Get(path string, defaultValue ...any) string {
	return GetString(path, defaultValue...)
}

func internalGet(path string, defaultValue ...any) any {
	// conf 或者环境变量不存在的情况
	if !v.IsSet(path) || v.Get(path) == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return nil
	}
	return v.Get(path)
}

// GetString 获取 String 类型的配置信息
func GetString(path string, defaultValue ...any) string {
	return cast.ToString(internalGet(path, defaultValue...))
}

// GetInt 获取 Int 类型的配置信息
func GetInt(path string, defaultValue ...any) int {
	return cast.ToInt(internalGet(path, defaultValue...))
}

// GetFloat64 获取 float64 类型的配置信息
func GetFloat64(path string, defaultValue ...any) float64 {
	return cast.ToFloat64(internalGet(path, defaultValue...))
}

// GetInt64 获取 Int64 类型的配置信息
func GetInt64(path string, defaultValue ...any) int64 {
	return cast.ToInt64(internalGet(path, defaultValue...))
}

// GetUint 获取 Uint 类型的配置信息
func GetUint(path string, defaultValue ...any) uint {
	return cast.ToUint(internalGet(path, defaultValue...))
}

// GetBool 获取 Bool 类型的配置信息
func GetBool(path string, defaultValue ...any) bool {
	return cast.ToBool(internalGet(path, defaultValue...))
}

// GetStringMapString 获取结构数据
func GetStringMapString(path string) map[string]string {
	return viper.GetStringMapString(path)
}
